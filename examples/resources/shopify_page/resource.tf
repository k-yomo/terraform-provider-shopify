resource "shopify_page" "example" {
  handle          = "example"
  author          = "Tom Brown"
  title           = "Example Page"
  body_html       = "<h1>Welcome to our store!</h1>"
  template_suffix = "page"
  published       = true
}
