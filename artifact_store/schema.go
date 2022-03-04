package artifact_store

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var ArtifactStore = map[string]*schema.Schema{
	"url": {
		Type:     schema.TypeString,
		Required: true,
	},
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
}
