package foxpass

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type providerConfiguration struct {
  Client                             *FoxpassClient
}

func Provider() *schema.Provider {
	return &schema.Provider{

    Schema: map[string]*schema.Schema{
			"api_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("FOXPASS_API_TOKEN", nil),
				Description: "API Token for Foxpass",
			},
    },

		ResourcesMap: map[string]*schema.Resource{
			"foxpass_mac_entry_prefix": MacEntryPrefix(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
    api_token := d.Get("api_token").(string)

    var diags diag.Diagnostics

	if (api_token == "") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Foxpass client",
			Detail:   "Foxpass api token not provided",
		})

		return nil, diags
	}

	return NewClient(api_token), diags
}
