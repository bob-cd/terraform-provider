package resource_provider

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataResourceProviders() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceProvidersRead,
		Schema: map[string]*schema.Schema{
			"resource_providers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: ResourceProvider,
				},
			},
		},
	}
}

func resourceProvidersRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	resourceProviders, err := GetAllResourceProviders()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := data.Set("resource_providers", resourceProviders); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
