package main

import (
	"context"
	"time"

	"github.com/bob-cd/terraform-provider/artifact_store"
	c "github.com/bob-cd/terraform-provider/common"
	"github.com/bob-cd/terraform-provider/pipeline"
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
			"reconcile_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BOB_RECONCILE_INTERVAL", 1000),
			},
			"reconcile_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BOB_RECONCILE_RETRIES", 10),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"bob_resource_provider": resource_provider.Resource(),
			"bob_artifact_store":    artifact_store.Resource(),
			"bob_pipeline":          pipeline.Resource(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, data *schema.ResourceData) (any, diag.Diagnostics) {
	url := data.Get("url").(string)
	timeout := time.Duration(data.Get("timeout").(int)) * time.Millisecond
	var diags diag.Diagnostics

	return c.NewClient(url, timeout, data.Get("reconcile_retries").(int), time.Duration(data.Get("reconcile_interval").(int))*time.Millisecond), diags
}
