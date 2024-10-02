// api/actions.go

package api

import (
	"bytes"
	"dcli/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// ListActions lists available actions filtered by entity type.
func (c *Client) ListActions(entityType string) ([]Action, error) {
	if entityType == "" {
		return nil, fmt.Errorf("entityType is required")
	}

	// Construct the URL with filters
	path := fmt.Sprintf("api/action?filter[OnType]=%s", url.QueryEscape(entityType))
	rel, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for key, values := range c.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list actions: %s\n%s", resp.Status, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON API response
	var jsonResponse struct {
		Data []struct {
			ID         string                 `json:"id"`
			Type       string                 `json:"type"`
			Attributes map[string]interface{} `json:"attributes"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse actions: %w", err)
	}

	var actions []Action
	for _, item := range jsonResponse.Data {
		var action Action
		action.ReferenceId = item.ID
		action.OnType = item.Type

		// Map attributes to Action struct
		err := mapToStruct(item.Attributes, &action)
		if err != nil {
			return nil, fmt.Errorf("failed to map action attributes: %w", err)
		}

		actions = append(actions, action)
	}

	return actions, nil
}

func (c *Client) GetAction(entityType, actionName string) (*Action, error) {
	if entityType == "" || actionName == "" {
		return nil, fmt.Errorf("entityType and actionName are required")
	}

	// Construct the URL with filters
	actionQuery, _ := json.Marshal([]map[string]interface{}{
		{
			"column":   "action_name",
			"value":    actionName,
			"operator": "eq",
		},
	})
	path := fmt.Sprintf("api/action?query=%s", actionQuery)
	rel, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for key, values := range c.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get action: %s\n%s", resp.Status, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON API response
	var jsonResponse struct {
		Data []models.Resource `json:"data"`
	}

	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse action: %w", err)
	}

	if len(jsonResponse.Data) == 0 {
		return nil, fmt.Errorf("action not found")
	}

	item := jsonResponse.Data[0]

	var action_schema = item.Attributes["action_schema"].(string)
	var action Action
	action.ReferenceId = item.ID
	err = json.Unmarshal([]byte(action_schema), &action)
	if err != nil {
		return nil, fmt.Errorf("failed to parse action schema: %w", err)
	}

	return &action, nil
}

func (c *Client) ExecuteAction(entityType, actionName string, inputs map[string]interface{}) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("action/%s/%s", url.PathEscape(entityType), url.PathEscape(actionName))
	rel, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	u := c.BaseURL.ResolveReference(rel)

	// Prepare request body
	data := inputs
	bodyData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(bodyData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add headers
	for key, values := range c.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("action execution failed: %s\n%s", resp.Status, string(body))
	}

	var result []map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse action execution result: %w", err)
	}

	return result, nil
}

// mapToStruct maps a map[string]interface{} to a struct
func mapToStruct(m map[string]interface{}, s interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, s)
}
