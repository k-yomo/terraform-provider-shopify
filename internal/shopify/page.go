package shopify

import goshopify "github.com/bold-commerce/go-shopify/v4"

func (c *Client) Page() goshopify.PageService {
	return c.shopifyClient.Page
}
