package aws

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

// tagsSchema returns the schema to use for tags.
//
func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
	}
}

// tagsSchema returns the schema to use for tags.
//
func tagsSchema2() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		ForceNew: true,
	}
}
func tagsSchemaComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}

// SetTagsDiff sets the new plan difference with the result of
// merging resource tags on to those defined at the provider-level;
// returns an error if unsuccessful or if the resource tags are identical
// to those configured at the provider-level to avoid non-empty plans
// after resource READ operations as resource and provider-level tags
// will be indistinguishable when returned from an AWS API.
func SetTagsDiff(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	defaultTagsConfig := meta.(*AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*AWSClient).IgnoreTagsConfig

	resourceTags := keyvaluetags.New(diff.Get("tags").(map[string]interface{}))

	if defaultTagsConfig.TagsEqual(resourceTags) {
		return fmt.Errorf(`"tags" are identical to those in the "default_tags" configuration block of the provider: please de-duplicate and try again`)
	}

	allTags := defaultTagsConfig.MergeTags(resourceTags).IgnoreConfig(ignoreTagsConfig)

	if err := diff.SetNew("tags_all", allTags.Map()); err != nil {
		return fmt.Errorf("error setting new tags_all diff: %w", err)
	}

	return nil
}
