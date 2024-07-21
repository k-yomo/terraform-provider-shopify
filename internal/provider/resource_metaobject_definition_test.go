package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMetaobjectDefinitionResource(t *testing.T) {
	metaobjectType := randResourceID(64)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccMetaobjectDefinitionResourceConfig(metaobjectType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("shopify_metaobject_definition.author", "name", "Author"),
					resource.TestCheckResourceAttr("shopify_metaobject_definition.author", "type", metaobjectType),
					resource.TestCheckResourceAttr("shopify_metaobject_definition.author", "field_definitions.#", "2"),
					resource.TestCheckResourceAttr("shopify_metaobject_definition.author", "field_definitions.0.key", "name"),
					resource.TestCheckResourceAttr("shopify_metaobject_definition.author", "field_definitions.0.name", "Name"),
					resource.TestCheckResourceAttr("shopify_metaobject_definition.author", "field_definitions.0.type", "single_line_text_field"),
					resource.TestCheckResourceAttr("shopify_metaobject_definition.author", "field_definitions.0.required", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "shopify_metaobject_definition.author",
				ImportState:       true,
				ImportStateVerify: true,
			},
			//// Update and Read testing
			{
				Config: testAccMetaobjectDefinitionResourceUpdateConfig(metaobjectType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("shopify_metaobject_definition.author", "name", "Updated Author"),
				),
			},
		},
	})
}

func testAccMetaobjectDefinitionResourceConfig(metaobjectType string) string {
	return fmt.Sprintf(`
resource "shopify_metaobject_definition" "author" {
  name       = "Author"
  type        = %[1]q
  field_definitions = [
    {
      key      = "name"
	  name     = "Name"
	  type     = "single_line_text_field"
      required = true
    },
	{
      key      = "profile_image_url"
	  name     = "Profile Image URL"
	  type     = "single_line_text_field"
    }
  ]
}
`, metaobjectType)
}

func testAccMetaobjectDefinitionResourceUpdateConfig(metaobjectType string) string {
	return fmt.Sprintf(`
resource "shopify_metaobject_definition" "author" {
  name       = "Updated Author"
  type        = %[1]q
  field_definitions = [
    {
      key      = "name"
	  name     = "Author Name"
	  type     = "single_line_text_field"
      required = true
    },
	{
      key      = "profile_image_url"
	  name     = "Profile Image URL"
	  type     = "url"
      required = true
    },
    {
      key      = "bio"
	  name     = "Bio"
	  type     = "rich_text_field"
    }
  ]
}
`, metaobjectType)
}
