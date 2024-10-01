// api/resource.go

package api

import (
	"encoding/json"
	"fmt"
	"jsonapi-cli-llm/models"
)

func (c *Client) Create(resource *models.Resource) (*models.Resource, error) {
	path := fmt.Sprintf("/%s", resource.Type)
	doc := &models.Document{
		Data: resource,
	}
	var respDoc models.Document
	err := c.post(path, doc, &respDoc)
	if err != nil {
		return nil, err
	}
	res, err := parseSingleResource(respDoc.Data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) Read(resourceType, id string) (*models.Resource, error) {
	path := fmt.Sprintf("/%s/%s", resourceType, id)
	var respDoc models.Document
	err := c.get(path, nil, &respDoc)
	if err != nil {
		return nil, err
	}
	res, err := parseSingleResource(respDoc.Data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) Update(resource *models.Resource) (*models.Resource, error) {
	if resource.ID == "" {
		return nil, fmt.Errorf("resource ID is required for update")
	}
	path := fmt.Sprintf("/%s/%s", resource.Type, resource.ID)
	doc := &models.Document{
		Data: resource,
	}
	var respDoc models.Document
	err := c.patch(path, doc, &respDoc)
	if err != nil {
		return nil, err
	}
	res, err := parseSingleResource(respDoc.Data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) Delete(resourceType, id string) error {
	path := fmt.Sprintf("/%s/%s", resourceType, id)
	err := c.delete(path)
	if err != nil {
		return err
	}
	return nil
}

func parseSingleResource(data interface{}) (*models.Resource, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var resource models.Resource
	err = json.Unmarshal(dataBytes, &resource)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}
