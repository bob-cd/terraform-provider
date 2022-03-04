package resource_provider

import (
	"context"

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

func read(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return c.ReadResource("resource-provider", data)
}

func create(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return c.CreateResource("resource-provider", data)
}

func update(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return c.UpdateResource("resource-provider", data)
}

func delete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return c.DeleteResource("resource-provider", data)
}
