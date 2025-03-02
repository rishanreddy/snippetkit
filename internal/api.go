package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// APIResponse represents the structure of the API response
type APIResponse struct {
	Success bool    `json:"success"`
	Data    []Snippet `json:"data"`
}

// Snippet represents the structure of a snippet returned from the API
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

// FetchSnippet fetches a snippet from the API
func FetchSnippet(snippetID string) (*Snippet, error) {
	url := fmt.Sprintf("http://localhost:3000/api/cli/snippet/get?id=%s", snippetID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response: %v", err)
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API response: %v", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API response was unsuccessful")
	}

	if len(apiResp.Data) == 0 || apiResp.Data[0].ID == "" {
		return nil, fmt.Errorf("no snippet found")
	}
	return &apiResp.Data[0], nil
}


// SearchSnippets searches for snippets by name, language, tags, and limit
func SearchSnippets(query, lang, tag string, limit int) ([]Snippet, error) {
	params := url.Values{}
	params.Add("q", query)
	if lang != "" {
		params.Add("lang", lang)
	}
	if tag != "" {
		params.Add("tag", tag)
	}
	params.Add("limit", fmt.Sprintf("%d", limit))

	apiURL := fmt.Sprintf("http://localhost:3000/api/cli/snippet/search?%s", params.Encode())
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response: %v", err)
	}

	var searchResp APIResponse
	err = json.Unmarshal(body, &searchResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API response: %v", err)
	}

	if !searchResp.Success {
		return nil, fmt.Errorf("API response was unsuccessful")
	}

	return searchResp.Data, nil
}