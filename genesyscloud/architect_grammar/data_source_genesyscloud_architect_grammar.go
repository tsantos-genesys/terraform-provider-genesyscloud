package architect_grammar

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	gcloud "terraform-provider-genesyscloud/genesyscloud"
)

/*
   The data_source_genesyscloud_grammar_proxy.go contains the data source implementation
   for the resource.
*/

// dataSourceArchitectGrammarRead retrieves by name the id in question
func dataSourceArchitectGrammarRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sdkConfig := meta.(*gcloud.ProviderMeta).ClientConfig
	proxy := newArchitectGrammarProxy(sdkConfig)

	name := d.Get("name").(string)

	return gcloud.WithRetries(ctx, 15*time.Second, func() *resource.RetryError {
		grammarId, retryable, err := proxy.getArchitectGrammarIdByName(ctx, name)

		if err != nil && !retryable {
			return resource.NonRetryableError(fmt.Errorf("Error grammar %s: %s", name, err))
		}

		if retryable {
			return resource.RetryableError(fmt.Errorf("No grammar found with name %s", name))
		}

		d.SetId(grammarId)
		return nil
	})
}
