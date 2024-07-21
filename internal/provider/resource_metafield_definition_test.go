package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMetafieldDefinitionResource(t *testing.T) {
	metafieldKey := randResourceID(64)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccMetafieldDefinitionResourceConfig(metafieldKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("shopify_metafield_definition.test", "key", metafieldKey),
					resource.TestCheckResourceAttr("shopify_metafield_definition.test", "name", "Terraform Test"),
					resource.TestCheckResourceAttr("shopify_metafield_definition.test", "namespace", "testacc"),
					resource.TestCheckResourceAttr("shopify_metafield_definition.test", "owner_type", "CUSTOMER"),
					resource.TestCheckResourceAttr("shopify_metafield_definition.test", "type", "single_line_text_field"),
					resource.TestCheckResourceAttr("shopify_metafield_definition.test", "pin", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "shopify_metafield_definition.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccMetafieldDefinitionResourceUpdateConfig(metafieldKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("shopify_metafield_definition.test", "name", "Terraform Test Updated"),
					resource.TestCheckResourceAttr("shopify_metafield_definition.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("shopify_metafield_definition.test", "pin", "true"),
					resource.TestCheckResourceAttr("shopify_metafield_definition.test", "validations.0.name", "min"),
					resource.TestCheckResourceAttr("shopify_metafield_definition.test", "validations.0.value", "10"),
				),
			},
		},
	})
}

func testAccMetafieldDefinitionResourceConfig(metafieldKey string) string {
	return fmt.Sprintf(`
resource "shopify_metafield_definition" "test" {
  key        = %[1]q
  name       = "Terraform Test"
  namespace  = "testacc"
  owner_type = "CUSTOMER"
  type       = "single_line_text_field"
}
`, metafieldKey)
}

func testAccMetafieldDefinitionResourceUpdateConfig(metafieldKey string) string {
	return fmt.Sprintf(`
resource "shopify_metafield_definition" "test" {
  key         = %[1]q
  name        = "Terraform Test Updated"
  description = "Updated description"
  namespace   = "testacc"
  owner_type  = "CUSTOMER"
  type        = "single_line_text_field"
  pin         = true
  validations = [
	{
	  name  = "min"
	  value = "10"	
    }
  ] 
}
`, metafieldKey)
}
