package oauth_client

import (
	"context"
	"fmt"
	"log"
	"terraform-provider-genesyscloud/genesyscloud/provider"
	"terraform-provider-genesyscloud/genesyscloud/util"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"terraform-provider-genesyscloud/genesyscloud/consistency_checker"
	"terraform-provider-genesyscloud/genesyscloud/util/resourcedata"

	resourceExporter "terraform-provider-genesyscloud/genesyscloud/resource_exporter"
	lists "terraform-provider-genesyscloud/genesyscloud/util/lists"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mypurecloud/platform-client-sdk-go/v125/platformclientv2"
)

func getAllOAuthClients(ctx context.Context, clientConfig *platformclientv2.Configuration) (resourceExporter.ResourceIDMetaMap, diag.Diagnostics) {
	resources := make(resourceExporter.ResourceIDMetaMap)

	oauthClientProxy := GetOAuthClientProxy(clientConfig)
	clients, resp, getErr := oauthClientProxy.getAllOAuthClients(ctx)

	if getErr != nil {
		return nil, diag.Errorf("Failed to get page of oauth clients: %v %v", getErr, resp)
	}

	for _, client := range *clients {
		if client.State != nil && *client.State == "disabled" {
			// Don't include clients disabled by support
			continue
		}
		resources[*client.Id] = &resourceExporter.ResourceMeta{Name: *client.Name}
	}
	return resources, nil
}

func createOAuthClient(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	tokenSeconds := d.Get("access_token_validity_seconds").(int)
	grantType := d.Get("authorized_grant_type").(string)
	state := d.Get("state").(string)

	sdkConfig := meta.(*provider.ProviderMeta).ClientConfig
	oauthClientProxy := GetOAuthClientProxy(sdkConfig)

	roles, diagErr := buildOAuthRoles(d)
	if diagErr != nil {
		return diagErr
	}

	//Before we create the oauth client take any roles on it and update the user
	diagErr = updateTerraformUserWithRole(ctx, sdkConfig, roles)
	if diagErr != nil {
		return diagErr
	}

	log.Printf("Creating oauth client %s", name)
	oauthRequest := &platformclientv2.Oauthclientrequest{
		Name:                       &name,
		Description:                &description,
		AccessTokenValiditySeconds: &tokenSeconds,
		AuthorizedGrantType:        &grantType,
		State:                      &state,
		RegisteredRedirectUri:      buildOAuthRedirectURIs(d),
		Scope:                      buildOAuthScopes(d),
		RoleDivisions:              roles,
	}

	client, resp, err := oauthClientProxy.createOAuthClient(ctx, *oauthRequest)
	if err != nil {
		return diag.Errorf("Failed to create oauth client %s: %s %v", name, err, resp)
	}

	credentialName := resourcedata.GetNillableValue[string](d, "integration_credential_name")
	if credentialName != nil {

		cred_type := "pureCloudOAuthClient"
		results := make(map[string]string)
		results["clientId"] = *client.Id
		results["clientSecret"] = *client.Secret

		createCredential := platformclientv2.Credential{
			Name: credentialName,
			VarType: &platformclientv2.Credentialtype{
				Name: &cred_type,
			},
			CredentialFields: &results,
		}

		credential, resp, err := oauthClientProxy.createIntegrationClient(ctx, createCredential)

		if err != nil {
			return diag.Errorf("Failed to create credential %s : %s %v", name, err, resp)
		}

		_ = d.Set("integration_credential_id", *credential.Id)
		_ = d.Set("integration_credential_name", *credential.Name)
	}

	d.SetId(*client.Id)
	log.Printf("Created oauth client %s %s", name, *client.Id)
	return readOAuthClient(ctx, d, meta)
}

func updateTerraformUserWithRole(ctx context.Context, sdkConfig *platformclientv2.Configuration, addedRoles *[]platformclientv2.Roledivision) diag.Diagnostics {
	op := GetOAuthClientProxy(sdkConfig)

	//Step #1 Retrieve the parent oauth client from the token API and check to make sure it is not a client credential grant
	tokenInfo, resp, err := op.getParentOAuthClientToken(ctx)
	if err != nil {
		return util.BuildAPIDiagnosticError(resourceName, fmt.Sprintf("Error trying to retrieve the token for the OAuth client running our CX as Code provider %s", err), resp)
	}

	if *tokenInfo.OAuthClient.Organization.Id != "purecloud-builtin" {
		log.Printf("This terraform client is being run with an OAuth Client Credential Grant.  You might get an error in your terraform scripts if you try to create a role in CX as Code and try to assign it to the oauth client.")
		return nil
	}

	//Step #2: Look up the user who is running the user
	log.Printf("The OAuth Client being used is purecloud-builtin. Retrieving the user running the terraform client and assigning the target role to them.")
	terraformUser, resp, err := op.GetTerraformUser(ctx)
	if err != nil {
		return util.BuildAPIDiagnosticError(resourceName, fmt.Sprintf("Failed to retrieved the terraform user running this terraform code %s", err), resp)
	}

	//Step #3: Lookup the users addedRoles
	userRoles, resp, err := op.GetTerraformUserRoles(ctx, *terraformUser.Id)
	if err != nil {
		return util.BuildAPIDiagnosticError(resourceName, fmt.Sprintf("Failued to retrieve the terraform user addedRoles running this terraform code %s", err), resp)
	}

	var totalRoles []string
	//Step #4  - Concat the addedRoles
	for _, role := range *addedRoles {
		totalRoles = append(totalRoles, *role.RoleId)
	}

	for _, role := range *userRoles.Roles {
		totalRoles = append(totalRoles, *role.Id)
	}

	//Step #5 - Update addedRoles
	_, resp, err = op.UpdateTerraformUserRoles(ctx, *terraformUser.Id, totalRoles)
	if err != nil {
		return util.BuildAPIDiagnosticError(resourceName, fmt.Sprintf("Failued to update the terraform user addedRoles running this terraform code %s", err), resp)
	}

	time.Sleep(5 * time.Second)
	return nil
}

func checkRoleAddedToUser(ctx context.Context, op *oauthClientProxy, terraformUser *platformclientv2.Userme, addedRoles *[]platformclientv2.Roledivision) diag.Diagnostics {
	log.Printf("Checking to see if the role has been added to the default user.")
	diagErr := util.WithRetries(ctx, 30*time.Second, func() *retry.RetryError {
		userAuth, _, err := op.GetTerraformUserRoles(ctx, *terraformUser.Id)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("Error while trying to validate if the user addedRoles have been updated "))
		}

		//Getting all the roles retrieved from the user and putting them into a map
		retrievedRolesMap := make(map[string]string)
		for _, role := range *userAuth.Roles {
			log.Printf("Adding a role to the retrievedRolesMap %s", *role.Id)
			retrievedRolesMap[*role.Id] = ""
		}

		//Going through each of the roles we added.  If the addedRole is not added, then return a retryable error
		for _, addedRole := range *addedRoles {
			if _, roleFound := retrievedRolesMap[*addedRole.RoleId]; !roleFound {
				log.Printf("Did ot match a role %s. Retrying", *addedRole.RoleId)
				return retry.RetryableError(fmt.Errorf("Was unable to find addedRole %s in the roles assigned to user.  Retrying read", *addedRole.RoleId))
			} else {
				log.Printf("Matched role: %s", *addedRole.RoleId)
			}

		}

		log.Printf("Succesfully matched all roles")
		return nil
	})
	return diagErr
}

func readOAuthClient(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sdkConfig := meta.(*provider.ProviderMeta).ClientConfig
	oAuthProxy := GetOAuthClientProxy(sdkConfig)

	log.Printf("Reading oauth client %s", d.Id())

	return util.WithRetriesForRead(ctx, d, func() *retry.RetryError {
		client, resp, getErr := oAuthProxy.getOAuthClient(ctx, d.Id())
		if getErr != nil {
			if util.IsStatus404(resp) {
				return retry.RetryableError(fmt.Errorf("Failed to read oauth client %s: %s", d.Id(), getErr))
			}
			return retry.NonRetryableError(fmt.Errorf("Failed to read oauth client %s: %s", d.Id(), getErr))
		}

		cc := consistency_checker.NewConsistencyCheck(ctx, d, meta, ResourceOAuthClient())
		_ = d.Set("name", *client.Name)

		resourcedata.SetNillableValue(d, "description", client.Description)
		resourcedata.SetNillableValue(d, "access_token_validity_seconds", client.AccessTokenValiditySeconds)
		resourcedata.SetNillableValue(d, "authorized_grant_type", client.AuthorizedGrantType)
		resourcedata.SetNillableValue(d, "state", client.State)

		if client.RegisteredRedirectUri != nil {
			_ = d.Set("registered_redirect_uris", lists.StringListToSet(*client.RegisteredRedirectUri))
		} else {
			_ = d.Set("registered_redirect_uris", nil)
		}

		if client.Scope != nil {
			_ = d.Set("scopes", lists.StringListToSet(*client.Scope))
		} else {
			_ = d.Set("scopes", nil)
		}

		if client.RoleDivisions != nil {
			_ = d.Set("roles", flattenOAuthRoles(*client.RoleDivisions))
		} else {
			_ = d.Set("roles", nil)
		}

		log.Printf("Read oauth client %s %s", d.Id(), *client.Name)
		return cc.CheckState()
	})
}

func updateOAuthClient(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	tokenSeconds := d.Get("access_token_validity_seconds").(int)
	grantType := d.Get("authorized_grant_type").(string)
	state := d.Get("state").(string)

	sdkConfig := meta.(*provider.ProviderMeta).ClientConfig
	oauthClientProxy := GetOAuthClientProxy(sdkConfig)

	roles, diagErr := buildOAuthRoles(d)
	if diagErr != nil {
		return diagErr
	}

	log.Printf("Updating oauth client %s", name)
	_, resp, err := oauthClientProxy.updateOAuthClient(ctx, d.Id(), platformclientv2.Oauthclientrequest{
		Name:                       &name,
		Description:                &description,
		AccessTokenValiditySeconds: &tokenSeconds,
		AuthorizedGrantType:        &grantType,
		State:                      &state,
		RegisteredRedirectUri:      buildOAuthRedirectURIs(d),
		Scope:                      buildOAuthScopes(d),
		RoleDivisions:              roles,
	})
	if err != nil {
		return diag.Errorf("Failed to update oauth client %s: %s %v", name, err, resp)
	}

	log.Printf("Updated oauth client %s", name)
	return readOAuthClient(ctx, d, meta)
}

func deleteOAuthClient(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sdkConfig := meta.(*provider.ProviderMeta).ClientConfig
	oauthClientProxy := GetOAuthClientProxy(sdkConfig)

	// check if there is a integration credential to delete
	credentialId := resourcedata.GetNillableValue[string](d, "integration_credential_id")
	if credentialId != nil {
		currentCredential, resp, getErr := oauthClientProxy.getIntegrationCredential(ctx, d.Id())
		if getErr == nil {
			_, err := oauthClientProxy.deleteIntegrationCredential(ctx, d.Id())
			return diag.Errorf("failed to delete integration credential %s (%s): %s %v", *currentCredential.Id, *currentCredential.Name, err, resp)
		}
	}

	name := d.Get("name").(string)

	log.Printf("Deleting oauth client %s", name)

	// The client state must be set to inactive before deleting
	_ = d.Set("state", "inactive")
	diagErr := updateOAuthClient(ctx, d, meta)
	if diagErr != nil {
		return diagErr
	}

	resp, err := oauthClientProxy.deleteOAuthClient(ctx, d.Id())
	if err != nil {
		return diag.Errorf("Failed to delete oauth client %s: %s %v", name, err, resp)
	}

	return util.WithRetries(ctx, 30*time.Second, func() *retry.RetryError {
		oauthClient, resp, err := oauthClientProxy.getOAuthClient(ctx, d.Id())
		if err != nil {
			if util.IsStatus404(resp) {
				// OAuth client deleted
				log.Printf("Deleted OAuth client %s", d.Id())
				return nil
			}
			return retry.NonRetryableError(fmt.Errorf("Error deleting OAuth client %s: %s", d.Id(), err))
		}

		if oauthClient.State != nil && *oauthClient.State == "deleted" {
			// OAuth client deleted
			log.Printf("Deleted OAuth client %s", d.Id())
			return nil
		}
		return retry.RetryableError(fmt.Errorf("OAuth client %s still exists", d.Id()))
	})
}
