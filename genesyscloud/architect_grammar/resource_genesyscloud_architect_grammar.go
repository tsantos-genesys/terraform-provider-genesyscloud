package architect_grammar

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mypurecloud/platform-client-sdk-go/v109/platformclientv2"
	"log"
	gcloud "terraform-provider-genesyscloud/genesyscloud"
	"terraform-provider-genesyscloud/genesyscloud/consistency_checker"
	resourceExporter "terraform-provider-genesyscloud/genesyscloud/resource_exporter"
	"terraform-provider-genesyscloud/genesyscloud/util/resourcedata"
)

/*
The resource_genesyscloud_architect_grammar.go contains all of the methods that perform the core logic for a resource.
*/

// getAllAuthArchitectGrammar retrieves all of the architect grammars via Terraform in the Genesys Cloud and is used for the exporter
func getAllAuthArchitectGrammar(ctx context.Context, clientConfig *platformclientv2.Configuration) (resourceExporter.ResourceIDMetaMap, diag.Diagnostics) {
	return nil, nil
}

// createArchitectGrammar is used by the architect_grammar resource to create a Genesys cloud architect grammar
func createArchitectGrammar(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sdkConfig := meta.(*gcloud.ProviderMeta).ClientConfig
	proxy := newArchitectGrammarProxy(sdkConfig)

	architectGrammar := getArchitectGrammarFromResourceData(d)

	grammar, err := proxy.createArchitectGrammar(ctx, &architectGrammar)
	if err != nil {
		return diag.Errorf("Failed to create grammar: %s", err)
	}

	d.SetId(*grammar.Id)
	log.Printf("Created Architect Grammar %s", *grammar.Id)
	return readArchitectGrammar(ctx, d, meta)
}

// readArchitectGrammar is used by the architect_grammar resource to read an architect grammar from genesys cloud.
func readArchitectGrammar(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sdkConfig := meta.(*gcloud.ProviderMeta).ClientConfig
	proxy := newArchitectGrammarProxy(sdkConfig)

	log.Printf("Reading Architect Grammar %s", d.Id())

	return gcloud.WithRetriesForRead(ctx, d, func() *resource.RetryError {
		grammar, respCode, getErr := proxy.getArchitectGrammarById(ctx, d.Id())
		if getErr != nil {
			if gcloud.IsStatus404ByInt(respCode) {
				return resource.RetryableError(fmt.Errorf("Failed to read Architect Grammar %s: %s", d.Id(), getErr))
			}
			return resource.NonRetryableError(fmt.Errorf("Failed to read Architect Grammar %s: %s", d.Id(), getErr))
		}

		cc := consistency_checker.NewConsistencyCheck(ctx, d, meta, ResourceArchitectGrammar())

		resourcedata.SetNillableValue(d, "name", grammar.Name)
		resourcedata.SetNillableValue(d, "description", grammar.Description)
		resourcedata.SetNillableValueWithInterfaceArrayWithFunc(d, "languages", grammar.Languages, flattenGrammarLanguages)

		log.Printf("Read Architect Grammar %s %s", d.Id(), *grammar.Name)
		return cc.CheckState()
	})
}

// updateArchitectGrammar is used by the architect_grammar resource to update an architect grammar in Genesys Cloud
func updateArchitectGrammar(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sdkConfig := meta.(*gcloud.ProviderMeta).ClientConfig
	proxy := newArchitectGrammarProxy(sdkConfig)

	architectGrammar := getArchitectGrammarFromResourceData(d)

	grammar, err := proxy.updateArchitectGrammar(ctx, d.Id(), &architectGrammar)
	if err != nil {
		return diag.Errorf("Failed to update grammar: %s", err)
	}

	log.Printf("Updated Architect Grammar %s", *grammar.Id)
	return readArchitectGrammar(ctx, d, meta)
}

// deleteArchitectGrammar is used by the architect_grammar resource to delete an architect grammar from Genesys cloud.
func deleteArchitectGrammar(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
