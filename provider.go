package main

import (
	"github.com/bob-cd/terraform-provider/artifact_store"
	"github.com/bob-cd/terraform-provider/resource_provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"bob_resource_provider": resource_provider.Resource(),
			"bob_artifact_store":    artifact_store.Resource(),
		},
	}
}
