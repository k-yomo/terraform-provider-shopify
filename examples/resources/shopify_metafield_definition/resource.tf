resource "shopify_metafield_definition" "example" {
  key        = "example"
  name       = "Example"
  namespace  = "custom"
  owner_type = "CUSTOMER"
  type       = "single_line_text_field"
}
