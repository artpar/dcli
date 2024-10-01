// models/article.go

package models

// ArticleAttributes defines the attributes for an article.
type ArticleAttributes struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Article represents an article resource.
type Article struct {
	Type          string                  `json:"type"`
	ID            string                  `json:"id,omitempty"`
	Attributes    ArticleAttributes       `json:"attributes,omitempty"`
	Relationships map[string]Relationship `json:"relationships,omitempty"`
	Links         *Links                  `json:"links,omitempty"`
	Meta          map[string]interface{}  `json:"meta,omitempty"`
}
