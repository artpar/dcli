// models/resource.go

package models

// Resource represents a JSON:API resource object.
type Resource struct {
	Type          string                  `json:"type"`
	ID            string                  `json:"id,omitempty"`
	Attributes    map[string]interface{}  `json:"attributes,omitempty"`
	Relationships map[string]Relationship `json:"relationships,omitempty"`
	Links         *Links                  `json:"links,omitempty"`
	Meta          map[string]interface{}  `json:"meta,omitempty"`
}

// ResourceIdentifier represents a resource identifier object.
type ResourceIdentifier struct {
	Type string                 `json:"type"`
	ID   string                 `json:"id"`
	Meta map[string]interface{} `json:"meta,omitempty"`
}

// Relationship represents a relationship object.
type Relationship struct {
	Links *Links                 `json:"links,omitempty"`
	Data  interface{}            `json:"data,omitempty"` // Can be ResourceIdentifier or []ResourceIdentifier
	Meta  map[string]interface{} `json:"meta,omitempty"`
}

// Links represents a links object.
type Links struct {
	Self    string `json:"self,omitempty"`
	Related string `json:"related,omitempty"`
}

// Document represents a JSON:API document.
type Document struct {
	Data     interface{}            `json:"data,omitempty"` // Can be Resource, []Resource, or ResourceIdentifier(s)
	Errors   []Error                `json:"errors,omitempty"`
	Meta     map[string]interface{} `json:"meta,omitempty"`
	JSONAPI  *JSONAPIObject         `json:"jsonapi,omitempty"`
	Links    *Links                 `json:"links,omitempty"`
	Included []Resource             `json:"included,omitempty"`
}

// JSONAPIObject represents a JSON:API object with version info.
type JSONAPIObject struct {
	Version string                 `json:"version,omitempty"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

// Error represents an error object as per JSON:API specification.
type Error struct {
	ID     string                 `json:"id,omitempty"`
	Links  *Links                 `json:"links,omitempty"`
	Status string                 `json:"status,omitempty"`
	Code   string                 `json:"code,omitempty"`
	Title  string                 `json:"title,omitempty"`
	Detail string                 `json:"detail,omitempty"`
	Source *ErrorSource           `json:"source,omitempty"`
	Meta   map[string]interface{} `json:"meta,omitempty"`
}

// ErrorSource represents the source of an error.
type ErrorSource struct {
	Pointer   string `json:"pointer,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}
