package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPageResource(t *testing.T) {
	pageHandle := randResourceID(64)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccPageResourceConfig(pageHandle),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("shopify_page.test", "handle", pageHandle),
					resource.TestCheckResourceAttr("shopify_page.test", "author", "Author"),
					resource.TestCheckResourceAttr("shopify_page.test", "title", "Test page"),
					resource.TestCheckResourceAttr("shopify_page.test", "body_html", "<h1>Test page</h1>"),
					resource.TestCheckResourceAttr("shopify_page.test", "template_suffix", ""),
					resource.TestCheckResourceAttr("shopify_page.test", "published", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName: "shopify_page.test",
				ImportState:  true,
			},
			//// Update and Read testing
			{
				Config: testAccPageResourceUpdateConfig(pageHandle),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("shopify_page.test", "handle", pageHandle),
					resource.TestCheckResourceAttr("shopify_page.test", "author", "Updated Author"),
					resource.TestCheckResourceAttr("shopify_page.test", "title", "Updated test page"),
					resource.TestCheckResourceAttr("shopify_page.test", "body_html", "<h1>Updated test page</h1>"),
					resource.TestCheckResourceAttr("shopify_page.test", "template_suffix", ""),
					resource.TestCheckResourceAttr("shopify_page.test", "published", "true"),
				),
			},
		},
	})
}

func testAccPageResourceConfig(pageHandle string) string {
	return fmt.Sprintf(`
resource "shopify_page" "test" {
  handle     = %[1]q
  author	 = "Author"
  title      = "Test page"
  body_html  = "<h1>Test page</h1>"
  template_suffix = ""
  published  = false
}
`, pageHandle)
}

func testAccPageResourceUpdateConfig(pageHandle string) string {
	return fmt.Sprintf(`
resource "shopify_page" "test" {
  handle     = %[1]q
  author	 = "Updated Author"
  title      = "Updated test page"
  body_html  = "<h1>Updated test page</h1>"
  template_suffix = ""
  published  = true
}
`, pageHandle)
}
