package resource_provider

import (
	"context"
	"time"

	c "github.com/bob-cd/terraform-provider/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   read,
		UpdateContext: update,
		DeleteContext: delete,
		Schema:        ResourceProvider,
	}
}

func read(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	providerName := data.Id()
	allProviders, err := c.FetchAll("resource-provider")
	if err != nil {
		return diag.FromErr(err)
	}

	for _, provider := range allProviders {
		if provider["name"] == providerName {
			data.Set("url", provider["url"])
			data.Set("name", provider["name"])
		}
	}

	return diags
}

func create(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	providerName := data.Get("name").(string)
	providerUrl := data.Get("url").(string)
	if err := c.Post("resource-provider", providerName, providerUrl); err != nil {
		return diag.FromErr(err)
	}

	if err := c.WaitForCondition(c.Reconcile("resource-provider", providerName, providerUrl), 10, 1*time.Second); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(providerName)

	read(ctx, data, m)

	return diags
}

func update(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	providerName := data.Get("name").(string)
	providerUrl := data.Get("url").(string)
	if err := c.Post("resource-provider", providerName, providerUrl); err != nil {
		return diag.FromErr(err)
	}

	if err := c.WaitForCondition(c.Reconcile("resource-provider", providerName, providerUrl), 10, 1*time.Second); err != nil {
		return diag.FromErr(err)
	}

	read(ctx, data, m)

	return diags
}

func delete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	if err := c.Delete("resource-provider", data.Get("name").(string)); err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")

	return diags
}
