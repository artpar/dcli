package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is the API client that performs all operations against the JSON:API server.
type Client struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
	Headers    http.Header
}

// NewClient creates a new API client with the specified base URL.
func NewClient(baseURL string) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	client := &Client{
		BaseURL: parsedURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Headers: make(http.Header),
	}

	// Set default headers
	client.Headers.Set("Content-Type", "application/vnd.api+json")
	client.Headers.Set("Accept", "application/vnd.api+json")

	return client, nil
}

// doRequest executes an HTTP request and decodes the response.
func (c *Client) doRequest(req *http.Request, v interface{}) error {
	// Add default headers to the request
	for key, values := range c.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		if v != nil {
			return json.NewDecoder(resp.Body).Decode(v)
		}
		return nil
	}

	// Handle error responses
	var apiError APIError
	if err := json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
		return fmt.Errorf("failed to decode error response: %w", err)
	}

	return errors.New(apiError.Error())
}

// APIError represents an error returned by the API.
type APIError struct {
	Errors []struct {
		Status string `json:"status"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
	} `json:"errors"`
}

// Error implements the error interface for APIError.
func (e *APIError) Error() string {
	if len(e.Errors) > 0 {
		return fmt.Sprintf("%s: %s", e.Errors[0].Title, e.Errors[0].Detail)
	}
	return "unknown API error"
}

// get sends a GET request.
func (c *Client) get(path string, queryParams map[string]string, v interface{}) error {
	rel, err := url.Parse(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	u := c.BaseURL.ResolveReference(rel)

	// Add query parameters
	q := u.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	return c.doRequest(req, v)
}

// post sends a POST request.
func (c *Client) post(path string, body interface{}, v interface{}) error {
	return c.sendRequestWithBody("POST", path, body, v)
}

// patch sends a PATCH request.
func (c *Client) patch(path string, body interface{}, v interface{}) error {
	return c.sendRequestWithBody("PATCH", path, body, v)
}

// delete sends a DELETE request.
func (c *Client) delete(path string) error {
	req, err := http.NewRequest("DELETE", c.BaseURL.ResolveReference(&url.URL{Path: path}).String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	return c.doRequest(req, nil)
}

// sendRequestWithBody sends a request with a JSON body.
func (c *Client) sendRequestWithBody(method, path string, body interface{}, v interface{}) error {
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

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	return c.doRequest(req, v)
}
