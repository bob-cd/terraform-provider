package pipeline

import (
	"context"
	"fmt"
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
		DeleteContext: deleteResource,
		Schema:        Pipeline,
	}
}

func get(m map[string]interface{}, key string, def interface{}) interface{} {
	value, exists := m[key]

	if !exists {
		return def
	}

	return value
}

func unwrapArtifactProduction(attrs map[string]interface{}) map[string]interface{} {
	steps := attrs["steps"].([]interface{})

	for _, step := range steps {
		s := step.(map[string]interface{})
		produces_artifact := s["produces_artifact"].([]interface{})

		if len(produces_artifact) == 0 {
			delete(s, "produces_artifact")
		} else {
			s["produces_artifact"] = produces_artifact[0]
		}
	}

	return attrs
}

func wrapArtifactProduction(attrs map[string]interface{}) map[string]interface{} {
	steps := attrs["steps"].([]interface{})

	for _, step := range steps {
		s := step.(map[string]interface{})

		val, exists := s["produces_artifact"]
		if exists {
			s["produces_artifact"] = []interface{}{val}
		} else {
			s["produces_artifact"] = []interface{}{}
		}
	}

	return attrs
}

func write(data *schema.ResourceData, client c.Client) error {
	group := data.Get("group").(string)
	name := data.Get("name").(string)
	attrs := map[string]interface{}{
		"image": data.Get("image"),
		"steps": data.Get("step"),
	}

	vars := data.Get("vars")
	if vars == nil {
		vars = map[string]string{}
	}
	attrs["vars"] = vars

	resources := data.Get("resource")
	if resources == nil {
		resources = []interface{}{}
	}
	attrs["resources"] = resources

	attrs = unwrapArtifactProduction(attrs) // Yes

	if err := client.PostPipeline(group, name, attrs); err != nil {
		return err
	}

	attrs["group"] = group
	attrs["name"] = name
	if err := c.WaitForCondition(client.ReconcilePipeline(attrs), 10, 1*time.Second); err != nil {
		return err
	}

	return nil
}

func read(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(c.Client)

	pipeline, err := client.FetchPipeline(data.Get("group").(string), data.Get("name").(string))
	if err != nil && err.Error() != "no such pipeline" {
		return diag.FromErr(err)
	}

	if err == nil {
		pipeline = wrapArtifactProduction(pipeline) // Yes

		data.Set("group", pipeline["group"])
		data.Set("name", pipeline["name"])
		data.Set("image", pipeline["image"])
		data.Set("vars", get(pipeline, "vars", map[string]string{}))
		data.Set("resource", get(pipeline, "resources", []interface{}{}))
		data.Set("step", pipeline["steps"])
	}

	return diags
}

func create(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := write(data, meta.(c.Client)); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(fmt.Sprintf("%s:%s", data.Get("group").(string), data.Get("name").(string)))

	read(ctx, data, meta)

	return diags
}

func update(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := write(data, meta.(c.Client)); err != nil {
		return diag.FromErr(err)
	}

	read(ctx, data, meta)

	return diags
}

func deleteResource(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(c.Client)
	group := data.Get("group").(string)
	name := data.Get("name").(string)

	if err := client.DeletePipeline(group, name); err != nil {
		return diag.FromErr(err)
	}

	if err := c.WaitForCondition(client.ReconcilePipelineDeletion(group, name), 10, 1*time.Second); err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")

	return diags
}
