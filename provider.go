package main

import (
	"context"
	"time"

	"github.com/bob-cd/terraform-provider/artifact_store"
	c "github.com/bob-cd/terraform-provider/common"
	"github.com/bob-cd/terraform-provider/resource_provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BOB_URL", "http://localhost:7777"),
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BOB_TIMEOUT", 10000),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"bob_resource_provider": resource_provider.Resource(),
			"bob_artifact_store":    artifact_store.Resource(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
	url := data.Get("url").(string)
	timeout := time.Duration(data.Get("timeout").(int)) * time.Millisecond
	var diags diag.Diagnostics

	return c.NewClient(url, timeout), diags
}
