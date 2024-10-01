// api/entity.go

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jsonapi-cli-llm/utils"
	"net/http"
	"net/url"
	"strings"
)

// ColumnInfo represents information about a column.
type ColumnInfo struct {
	Name              string         `json:"Name"`
	ColumnName        string         `json:"ColumnName"`
	ColumnDescription string         `json:"ColumnDescription"`
	ColumnType        string         `json:"ColumnType"`
	DataType          string         `json:"DataType"`
	ForeignKeyData    ForeignKeyData `json:"ForeignKeyData"` // Add this line
	// Include other fields as necessary
}

// ForeignKeyData represents foreign key information for a column.
type ForeignKeyData struct {
	DataSource string `json:"DataSource"`
	Namespace  string `json:"Namespace"`
	KeyName    string `json:"KeyName"`
}

// Action represents an action that can be performed on the entity.
type Action struct {
	Name        string `json:"Name"`
	Label       string `json:"Label"`
	Description string `json:"Description"`
	// Add other fields as necessary
}

// Relationship represents a relationship with another entity.
type Relationship struct {
	RelationType string `json:"RelationType"` // e.g., hasOne, hasMany
	RelatedType  string `json:"RelatedType"`  // e.g., user_account
	// Add other fields as necessary
}

// EntityModel represents the model of an entity.
type EntityModel struct {
	ColumnModel   map[string]ColumnInfo `json:"ColumnModel"`
	Actions       []Action              `json:"Actions"`
	Relationships map[string]Relationship
	// Add other fields as necessary
}

// GetEntityModel fetches the entity model from the server.
func (c *Client) GetEntityModel(entityName string) (*EntityModel, error) {
	path := fmt.Sprintf("jsmodel/%s.js", entityName)
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

	// Optional: Log the request URL
	utils.InfoLogger.Printf("Request URL: %s", req.URL)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("failed to fetch entity model: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var model EntityModel
	err = json.Unmarshal(body, &model)
	if err != nil {
		return nil, fmt.Errorf("failed to parse entity model: %w", err)
	}

	// Process relationships
	model.Relationships = make(map[string]Relationship)
	for name, col := range model.ColumnModel {
		if col.ColumnType == "entity" || col.ColumnType == "alias" {
			relationType := "hasOne"
			if col.ForeignKeyData.Namespace != "" {
				relatedType := col.ForeignKeyData.Namespace
				if col.ColumnType == "entity" && strings.Contains(col.DataType, "[]") {
					relationType = "hasMany"
				}
				model.Relationships[name] = Relationship{
					RelationType: relationType,
					RelatedType:  relatedType,
				}
			}
		}
	}

	return &model, nil
}
