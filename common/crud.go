package common

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func write(entity string, data *schema.ResourceData, c Client) error {
	entityName := data.Get("name").(string)
	entityUrl := data.Get("url").(string)
	if err := c.Post(entity, entityName, entityUrl); err != nil {
		return err
	}

	if err := WaitForCondition(c.Reconcile(entity, entityName, entityUrl), c.ReconcileRetries, c.ReconcileInterval); err != nil {
		return err
	}

	return nil
}

func ReadResource(entity string, data *schema.ResourceData, c Client) diag.Diagnostics {
	var diags diag.Diagnostics

	entityName := data.Id()
	allEntities, err := c.FetchAll(entity)
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

func CreateResource(entity string, data *schema.ResourceData, c Client) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := write(entity, data, c); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(data.Get("name").(string))

	ReadResource(entity, data, c)

	return diags
}

func UpdateResource(entity string, data *schema.ResourceData, c Client) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := write(entity, data, c); err != nil {
		return diag.FromErr(err)
	}

	ReadResource(entity, data, c)

	return diags
}

func DeleteResource(entity string, data *schema.ResourceData, c Client) diag.Diagnostics {
	var diags diag.Diagnostics
	entityName := data.Get("name").(string)
	entityUrl := data.Get("url").(string)

	if err := c.Delete(entity, entityName); err != nil {
		return diag.FromErr(err)
	}

	if err := WaitForCondition(Complement(c.Reconcile(entity, entityName, entityUrl)), c.ReconcileRetries, c.ReconcileInterval); err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")

	return diags
}
