package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// APIResponseSingle is used for FetchSnippet (returns a single snippet)
type APIResponseSingle struct {
	Success bool    `json:"success"`
	Data    Snippet `json:"data"` // Expecting a single object
}

// APIResponseMultiple is used for SearchSnippets (returns multiple snippets)
type APIResponseMultiple struct {
	Success bool      `json:"success"`
	Data    []Snippet `json:"data"` // Expecting an array
}

// Snippet represents a single snippet
type Snippet struct {
	ID          string   `json:"_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Language    string   `json:"language"`
	ShortID     string   `json:"shortId"`
	Code        string   `json:"code"`
	Path        string   `json:"path"`
	Tags        []string `json:"tags"`
}

// FetchSnippet fetches a single snippet from the API
func FetchSnippet(snippetID string, apiKey string) (*Snippet, error) {
	apiURL := fmt.Sprintf("http://localhost:3000/api/snippet/get/%s", snippetID)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response: %v", err)
	}

	var apiResp APIResponseSingle
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API response: %v", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API response was unsuccessful")
	}

	return &apiResp.Data, nil
}

// SearchSnippets searches for snippets by query, language, tag, and limit
func SearchSnippets(query, lang, tag string, limit int, apiKey string) ([]Snippet, error) {
	// Build query parameters
	params := url.Values{}
	params.Add("q", query)
	if lang != "" {
		params.Add("lang", lang)
	}
	if tag != "" {
		params.Add("tag", tag)
	}
	params.Add("limit", fmt.Sprintf("%d", limit))

	apiURL := fmt.Sprintf("http://localhost:3000/api/snippet/search?%s", params.Encode())

	// Create an HTTP request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response: %v", err)
	}

	var searchResp APIResponseMultiple // Expecting an array in "data"
	err = json.Unmarshal(body, &searchResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API response: %v", err)
	}

	if !searchResp.Success {
		return nil, fmt.Errorf("API response was unsuccessful")
	}

	return searchResp.Data, nil
}

// VerifyToken verifies the API key by calling the `/api/token/verify` endpoint
func VerifyToken(apiKey string) (bool, error) {
	apiURL := "http://localhost:3000/api/token/verify"

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %v", err)
	}

	// Set API key in headers
	req.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to connect to API: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read API response: %v", err)
	}

	// Parse the JSON response
	var apiResp struct {
		Success bool   `json:"success"`
		Error   string `json:"error,omitempty"`
	}
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return false, fmt.Errorf("failed to parse API response: %v", err)
	}

	// If success is true, return true
	if apiResp.Success {
		return true, nil
	}

	// If not successful, return the error message
	return false, fmt.Errorf("API key validation failed: %s", apiResp.Error)
}