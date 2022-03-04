package resource_provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

func doPost(name string, url string) error {
	postBody, _ := json.Marshal(map[string]string{
		"url": url,
	})

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/resource-providers/%s", Host, name),
		bytes.NewBuffer(postBody),
	)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return err
	}

	_, err = Client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func reconcile(name string, url string) func() bool {
	return func() bool {
		allProviders, err := GetAllResourceProviders()
		if err != nil {
			return false
		}

		for _, provider := range allProviders {
			if provider["name"] == name && provider["url"] == url {
				return true
			}
		}

		return false
	}
}

func read(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	providerName := data.Id()
	allProviders, err := GetAllResourceProviders()
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
	if err := doPost(providerName, providerUrl); err != nil {
		return diag.FromErr(err)
	}

	if err := WaitForCondition(reconcile(providerName, providerUrl), 10, 1*time.Second); err != nil {
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
	if err := doPost(providerName, providerUrl); err != nil {
		return diag.FromErr(err)
	}

	if err := WaitForCondition(reconcile(providerName, providerUrl), 10, 1*time.Second); err != nil {
		return diag.FromErr(err)
	}

	read(ctx, data, m)

	return diags
}

func delete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/resource-providers/%s", Host, data.Get("name").(string)),
		nil,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = Client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")

	return diags
}
