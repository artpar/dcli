// api/entity.go

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jsonapi-cli-llm/utils"
	"net/http"
	"net/url"
)

// ForeignKeyData represents foreign key information for a column.
type ForeignKeyData struct {
	DataSource string `json:"DataSource"`
	Namespace  string `json:"Namespace"`
	KeyName    string `json:"KeyName"`
}

type Action struct {
	Name                    string       `json:"Name"`
	Label                   string       `json:"Label"`
	Description             string       `json:"Description"`
	OnType                  string       `json:"OnType"`
	InstanceOptional        bool         `json:"InstanceOptional"`
	RequestSubjectRelations []string     `json:"RequestSubjectRelations"`
	ReferenceId             string       `json:"ReferenceId"`
	InFields                []ColumnInfo `json:"InFields"`
	OutFields               []Outcome    `json:"OutFields"`
	Validations             []ColumnTag  `json:"Validations"`
	Conformations           []ColumnTag  `json:"Conformations"`
	// Add other fields as necessary
}

// Outcome and ColumnTag are placeholders; define them as needed

type Outcome struct {
	Type            string
	Method          string // method name
	Reference       string
	SkipInResponse  bool
	Condition       string
	Attributes      map[string]interface{}
	ContinueOnError bool
}

type ColumnTag struct {
	ColumnName string `json:"ColumnName"`
	Tags       string `json:"Tags"`
}

// Relationship represents a relationship with another entity.
type Relationship struct {
	RelationType string `json:"RelationType"` // e.g., hasOne, hasMany
	RelatedType  string `json:"RelatedType"`  // e.g., user_account
	// Add other fields as necessary
}
type AuthPermission int64

// TableInfo represents the model of an entity.

type TableInfo struct {
	TableName              string              `db:"table_name" json:"TableName"`
	TableId                int                 `json:"TableId"`
	DefaultPermission      AuthPermission      `db:"default_permission" json:"DefaultPermission"`
	Columns                []ColumnInfo        `json:"Columns"`
	Relations              []TableRelation     `json:"Relations"`
	IsTopLevel             bool                `db:"is_top_level" json:"IsTopLevel"`
	Permission             AuthPermission      `json:"Permission"`
	UserId                 uint64              `db:"user_account_id" json:"UserId"`
	IsHidden               bool                `db:"is_hidden" json:"IsHidden"`
	IsJoinTable            bool                `db:"is_join_table" json:"IsJoinTable"`
	IsStateTrackingEnabled bool                `db:"is_state_tracking_enabled" json:"IsStateTrackingEnabled"`
	IsAuditEnabled         bool                `db:"is_audit_enabled" json:"IsAuditEnabled"`
	TranslationsEnabled    bool                `db:"translation_enabled" json:"TranslationsEnabled"`
	DefaultGroups          []string            `db:"default_groups" json:"DefaultGroups"`
	DefaultRelations       map[string][]string `db:"default_relations" json:"DefaultRelations"`
	Validations            []ColumnTag         `json:"Validations"`
	Conformations          []ColumnTag         `json:"Conformations"`
	DefaultOrder           string              `json:"DefaultOrder"`
	Icon                   string              `json:"Icon"`
	CompositeKeys          [][]string          `json:"CompositeKeys"`
	Actions                []Action            `json:"Actions"` // Add this line
}

type ColumnInfo struct {
	Name              string         `db:"name" json:"Name"`
	ColumnName        string         `db:"column_name" json:"ColumnName"`
	ColumnDescription string         `db:"column_description" json:"ColumnDescription"`
	ColumnType        string         `db:"column_type" json:"ColumnType"`
	IsPrimaryKey      bool           `db:"is_primary_key" json:"IsPrimaryKey"`
	IsAutoIncrement   bool           `db:"is_auto_increment" json:"IsAutoIncrement"`
	IsIndexed         bool           `db:"is_indexed" json:"IsIndexed"`
	IsUnique          bool           `db:"is_unique" json:"IsUnique"`
	IsNullable        bool           `db:"is_nullable" json:"IsNullable"`
	Permission        uint64         `db:"permission" json:"Permission"`
	IsForeignKey      bool           `db:"is_foreign_key" json:"IsForeignKey"`
	ExcludeFromApi    bool           `db:"exclude_from_api" json:"ExcludeFromApi"`
	ForeignKeyData    ForeignKeyData `db:"foreign_key_data" json:"ForeignKeyData"`
	DataType          string         `db:"data_type" json:"DataType"`
	DefaultValue      string         `db:"default_value" json:"DefaultValue"`
	Options           []ValueOptions `json:"Options"`
	JsonApi           string         `json:"jsonApi"` // For relations
	Type              string         `json:"type"`    // Related entity type
}

type ValueOptions struct {
	ValueType string
	Value     interface{}
	Label     string
}

type TableRelation struct {
	Subject     string
	Object      string
	Relation    string
	SubjectName string
	ObjectName  string
	Columns     []ColumnInfo
}

// GetEntityModel fetches the entity model from the server.
func (c *Client) GetEntityModel(entityName string) (*TableInfo, error) {
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

	var model TableInfo
	err = json.Unmarshal(body, &model)
	if err != nil {
		return nil, fmt.Errorf("failed to parse entity model: %w", err)
	}

	return &model, nil
}
