resource "shopify_metaobject_definition" "example" {
  name = "Example"
  type = "example"
  field_definitions = [
    {
      key      = "text_field"
      name     = "Text Field"
      type     = "single_line_text_field"
      required = true
    },
    {
      key      = "url_field"
      name     = "URL Field"
      type     = "url"
      required = true
    },
    {
      key  = "rich_text_field"
      name = "Rich Text Field"
      type = "rich_text_field"
    }
  ]
}
