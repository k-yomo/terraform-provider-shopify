---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "shopify_page Resource - terraform-provider-shopify"
subcategory: ""
description: |-
  Page definitions enable you to define additional validation constraints for metafields, and enable the merchant to edit metafield values in context.
---

# shopify_page (Resource)

Page definitions enable you to define additional validation constraints for metafields, and enable the merchant to edit metafield values in context.

## Example Usage

```terraform
resource "shopify_page" "example" {
  handle          = "example"
  author          = "Tom Brown"
  title           = "Example Page"
  body_html       = "<h1>Welcome to our store!</h1>"
  template_suffix = "page"
  published       = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `author` (String) The name of the person who created the page.
- `body_html` (String) The text content of the page, complete with HTML markup.
- `handle` (String) A unique, human-friendly string for the page, generated automatically from its title. In themes, the Liquid templating language refers to a page by its handle.
- `title` (String) The title of the page.

### Optional

- `published` (Boolean) Whether the page is published. If true, the page is visible to customers. If false, the page is hidden from customers.
- `template_suffix` (String) he suffix of the template that is used to render the page. If the value is an empty string or null, then the default page template is used.

### Read-Only

- `id` (String) The unique numeric identifier for the page.
- `published_at` (String) The date and time (ISO 8601 format) when the page was published.

## Import

Import is supported using the following syntax:

```shell
# Note: integer id instead of graphql global id
terraform import shopify_page.test {{id}}
```
