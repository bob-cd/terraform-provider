package pipeline

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var Pipeline = map[string]*schema.Schema{
	"group": {
		Type:     schema.TypeString,
		Required: true,
	},
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"image": {
		Type:     schema.TypeString,
		Required: true,
	},
	"vars": {
		Type:     schema.TypeMap,
		Optional: true,
	},
	"resource": { // Singular as its gonna be series of blocks in HCL
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"type": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"internal", "external"}, false),
				},
				"provider": {
					Type:     schema.TypeString,
					Required: true,
				},
				"params": {
					Type:     schema.TypeMap,
					Optional: true,
					Default:  map[string]string{},
				},
			},
		},
	},
	"step": { // Singular as its gonna be series of blocks in HCL
		Type:     schema.TypeList,
		Required: true,
		MinItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"cmd": {
					Type:     schema.TypeString,
					Required: true,
				},
				"needs_resource": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"produces_artifact": {
					Type:     schema.TypeList,
					Optional: true,
          MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:     schema.TypeString,
								Required: true,
							},
							"path": {
								Type:     schema.TypeString,
								Required: true,
							},
							"store": {
								Type:     schema.TypeString,
								Required: true,
							},
						},
					},
				},
			},
		},
	},
}
