// api/pagination.go

package api

import (
	"fmt"
	"jsonapi-cli-llm/models"
)

type ListOptions struct {
	Page    map[string]string
	Filter  map[string]string
	Sort    string
	Include string
	Fields  map[string]string
}

func (c *Client) List(resourceType string, options *ListOptions) (*models.Document, error) {
	path := fmt.Sprintf("api/%s", resourceType)
	queryParams := make(map[string]string)

	if options != nil {
		for k, v := range options.Page {
			queryParams[fmt.Sprintf("page[%s]", k)] = v
		}
		for k, v := range options.Filter {
			queryParams[fmt.Sprintf("filter[%s]", k)] = v
		}
		if options.Sort != "" {
			queryParams["sort"] = options.Sort
		}
		if options.Include != "" {
			queryParams["include"] = options.Include
		}
		for k, v := range options.Fields {
			queryParams[fmt.Sprintf("fields[%s]", k)] = v
		}
	}

	var respDoc models.Document
	err := c.get(path, queryParams, &respDoc)
	if err != nil {
		return nil, err
	}
	return &respDoc, nil
}
