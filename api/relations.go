// api/relations.go

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jsonapi-cli-llm/models"
	"net/http"
	"net/url"
)

func (c *Client) FetchRelations(resourceType, id, relation string) (*models.Document, error) {
	path := fmt.Sprintf("%s/%s/%s", resourceType, id, relation)
	var respDoc models.Document
	err := c.get(path, nil, &respDoc)
	if err != nil {
		return nil, err
	}
	return &respDoc, nil
}

func (c *Client) GetRelationship(resourceType, id, relation string) (*models.Document, error) {
	path := fmt.Sprintf("%s/%s/relationships/%s", resourceType, id, relation)
	var respDoc models.Document
	err := c.get(path, nil, &respDoc)
	if err != nil {
		return nil, err
	}
	return &respDoc, nil
}

func (c *Client) UpdateRelationship(resourceType, id, relation string, data interface{}) (*models.Document, error) {
	path := fmt.Sprintf("%s/%s/relationships/%s", resourceType, id, relation)
	doc := &models.Document{
		Data: data,
	}
	var respDoc models.Document
	err := c.patch(path, doc, &respDoc)
	if err != nil {
		return nil, err
	}
	return &respDoc, nil
}

func (c *Client) AddToRelationship(resourceType, id, relation string, data interface{}) (*models.Document, error) {
	path := fmt.Sprintf("%s/%s/relationships/%s", resourceType, id, relation)
	doc := &models.Document{
		Data: data,
	}
	var respDoc models.Document
	err := c.post(path, doc, &respDoc)
	if err != nil {
		return nil, err
	}
	return &respDoc, nil
}

func (c *Client) DeleteFromRelationship(resourceType, id, relation string, data interface{}) error {
	path := fmt.Sprintf("%s/%s/relationships/%s", resourceType, id, relation)
	doc := &models.Document{
		Data: data,
	}
	err := c.deleteWithBody(path, doc)
	if err != nil {
		return err
	}
	return nil
}

// In api/client.go, add the deleteWithBody method

// deleteWithBody sends a DELETE request with a body.
func (c *Client) deleteWithBody(path string, body interface{}) error {
	rel, err := url.Parse(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	u := c.BaseURL.ResolveReference(rel)

	var buf io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		buf = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest("DELETE", u.String(), buf)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	return c.doRequest(req, nil)
}
