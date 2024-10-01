// api/filtering.go

package api

import (
	"dcli/models"
	"fmt"
)

func (c *Client) Filter(resourceType string, filters map[string]string) (*models.Document, error) {
	path := fmt.Sprintf("api/%s", resourceType)
	queryParams := make(map[string]string)

	for k, v := range filters {
		queryParams[fmt.Sprintf("filter[%s]", k)] = v
	}

	var respDoc models.Document
	err := c.get(path, queryParams, &respDoc)
	if err != nil {
		return nil, err
	}
	return &respDoc, nil
}
