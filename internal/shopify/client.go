package shopify

import (
	goshopify "github.com/bold-commerce/go-shopify/v4"
)

type Client struct {
	shopifyClient *goshopify.Client
}

func NewClient(shopifyClient *goshopify.Client) *Client {
	return &Client{
		shopifyClient: shopifyClient,
	}
}
