package shopify

import (
	"context"
)

type MetafieldDefinition struct {
	ID             string                   `json:"id"`
	Name           string                   `json:"name"`
	Description    *string                  `json:"description"`
	OwnerType      string                   `json:"ownerType"`
	Namespace      string                   `json:"namespace"`
	Key            string                   `json:"key"`
	Type           *MetafieldDefinitionType `json:"type"`
	PinnedPosition *int                     `json:"pinnedPosition"`
}

type MetafieldDefinitionType struct {
	Category string `json:"category"`
	Name     string `json:"name"`
}

type CreateMetafieldDefinitionInput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	OwnerType   string `json:"ownerType"`
	Namespace   string `json:"namespace"`
	Key         string `json:"key"`
	Type        string `json:"type"`
	Pin         bool   `json:"pin"`
}

type CreateMetafieldDefinitionResponse struct {
	MetafieldDefinitionCreate struct {
		CreatedDefinition *MetafieldDefinition `json:"createdDefinition"`
		UserErrors        UserErrors           `json:"userErrors"`
	} `json:"metafieldDefinitionCreate"`
}

func (c *Client) CreateMetafieldDefinition(ctx context.Context, input *CreateMetafieldDefinitionInput) (*MetafieldDefinition, error) {
	variables := map[string]interface{}{"definition": input}
	query := `
mutation CreateMetafieldDefinition($definition: MetafieldDefinitionInput!) {
  metafieldDefinitionCreate(definition: $definition) {
    createdDefinition {
      id
      name
      description
      ownerType
      namespace
      key
      type {
        category
        name
      }
      pinnedPosition
    }
    userErrors {
      field
      message
      code
    }
  }
}`

	var gqlResp CreateMetafieldDefinitionResponse
	err := c.shopifyClient.GraphQL.Query(ctx, query, variables, &gqlResp)
	if err != nil {
		return nil, err
	}
	if err := gqlResp.MetafieldDefinitionCreate.UserErrors.Error(); err != nil {
		return nil, err
	}
	return gqlResp.MetafieldDefinitionCreate.CreatedDefinition, nil
}

type GetMetafieldDefinitionResponse struct {
	MetafieldDefinition *MetafieldDefinition `json:"metafieldDefinition"`
}

func (c *Client) GetMetafieldDefinition(ctx context.Context, id string) (*MetafieldDefinition, error) {
	variables := map[string]interface{}{"id": id}
	query := `
query metafieldDefinition($id: ID!) {
  metafieldDefinition(id: $id) {
	id
    name
	description
	key
	namespace
	ownerType
    type {
      category
      name
    }
	pinnedPosition
  }
}
`

	var gqlResp GetMetafieldDefinitionResponse
	err := c.shopifyClient.GraphQL.Query(ctx, query, variables, &gqlResp)
	if err != nil {
		return nil, err
	}
	return gqlResp.MetafieldDefinition, nil
}

type UpdateMetafieldDefinitionInput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	OwnerType   string `json:"ownerType"`
	Namespace   string `json:"namespace"`
	Key         string `json:"key"`
	Pin         bool   `json:"pin"`
}

type UpdateMetafieldDefinitionResponse struct {
	MetafieldDefinitionUpdate struct {
		UpdatedDefinition *MetafieldDefinition `json:"updatedDefinition"`
		UserErrors        UserErrors           `json:"userErrors"`
	} `json:"metafieldDefinitionUpdate"`
}

func (c *Client) UpdateMetafieldDefinition(ctx context.Context, input *UpdateMetafieldDefinitionInput) (*MetafieldDefinition, error) {
	variables := map[string]interface{}{"definition": input}
	query := `
mutation UpdateMetafieldDefinition($definition: MetafieldDefinitionUpdateInput!) {
  metafieldDefinitionUpdate(definition: $definition) {
    updatedDefinition {
      id
      name
      description
      ownerType
      namespace
      key
      type {
        category
        name
      }
      pinnedPosition
    }
    userErrors {
      field
      message
      code
    }
  }
}`

	var gqlResp UpdateMetafieldDefinitionResponse
	err := c.shopifyClient.GraphQL.Query(ctx, query, variables, &gqlResp)
	if err != nil {
		return nil, err
	}
	if err := gqlResp.MetafieldDefinitionUpdate.UserErrors.Error(); err != nil {
		return nil, err
	}
	return gqlResp.MetafieldDefinitionUpdate.UpdatedDefinition, nil
}

type DeleteMetafieldDefinitionResponse struct {
	MetafieldDefinitionDelete struct {
		DeletedDefinitionID string `json:"DeletedDefinitionId"`
		UserErrors          UserErrors
	} `json:"metafieldDefinitionDelete"`
}

func (c *Client) DeleteMetafieldDefinition(ctx context.Context, id string) error {
	variables := map[string]interface{}{"id": id}
	query := `
mutation DeleteMetafieldDefinition($id: ID!) {
  metafieldDefinitionDelete(id: $id) {
	deletedDefinitionId
    userErrors {
      field
      message
      code
    }
  }
}`

	var gqlResp DeleteMetafieldDefinitionResponse
	err := c.shopifyClient.GraphQL.Query(ctx, query, variables, &gqlResp)
	if err != nil {
		return err
	}
	if err := gqlResp.MetafieldDefinitionDelete.UserErrors.Error(); err != nil {
		return err
	}
	return nil
}
