package tavily

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

type SearchRequest struct {
	Query          string   `json:"query"`
	SearchDepth    string   `json:"search_depth,omitempty"` // "basic" o "advanced"
	MaxResults     int      `json:"max_results,omitempty"`
	IncludeDomains []string `json:"include_domains,omitempty"`
	ExcludeDomains []string `json:"exclude_domains,omitempty"`
}

type SearchResult struct {
	Title   string  `json:"title"`
	URL     string  `json:"url"`
	Content string  `json:"content"`
	Score   float64 `json:"score"`
}

type SearchResponse struct {
	Results []SearchResult `json:"results"`
	Query   string         `json:"query"`
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: "https://api.tavily.com",
		client:  &http.Client{},
	}
}

func (c *Client) Search(req SearchRequest) (*SearchResponse, error) {
	// Valores por defecto
	if req.SearchDepth == "" {
		req.SearchDepth = "basic"
	}
	if req.MaxResults == 0 {
		req.MaxResults = 5
	}

	payload := map[string]interface{}{
		"api_key":      c.apiKey,
		"query":        req.Query,
		"search_depth": req.SearchDepth,
		"max_results":  req.MaxResults,
	}

	if len(req.IncludeDomains) > 0 {
		payload["include_domains"] = req.IncludeDomains
	}
	if len(req.ExcludeDomains) > 0 {
		payload["exclude_domains"] = req.ExcludeDomains
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/search", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var searchResp SearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return &searchResp, nil
}
