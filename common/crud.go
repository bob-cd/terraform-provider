package common

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ReadResource(entity string, data *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	entityName := data.Id()
	allEntities, err := FetchAll(entity)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, entity := range allEntities {
		if entity["name"] == entityName {
			data.Set("url", entity["url"])
			data.Set("name", entity["name"])
		}
	}

	return diags
}

func CreateResource(entity string, data *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	entityName := data.Get("name").(string)
	entityUrl := data.Get("url").(string)
	if err := Post(entity, entityName, entityUrl); err != nil {
		return diag.FromErr(err)
	}

	if err := WaitForCondition(Reconcile(entity, entityName, entityUrl), 10, 1*time.Second); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(entityName)

	ReadResource(entity, data)

	return diags
}

func UpdateResource(entity string, data *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	entityName := data.Get("name").(string)
	entityUrl := data.Get("url").(string)
	if err := Post("resource-provider", entityName, entityUrl); err != nil {
		return diag.FromErr(err)
	}

	if err := WaitForCondition(Reconcile(entity, entityName, entityUrl), 10, 1*time.Second); err != nil {
		return diag.FromErr(err)
	}

	ReadResource(entity, data)

	return diags
}

func DeleteResource(entity string, data *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	if err := Delete(entity, data.Get("name").(string)); err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")

	return diags
}
