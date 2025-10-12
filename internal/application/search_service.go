package application

import "github.com/Maxito7/hotel_backend/internal/tavily"

type SearchService struct {
	tavilyClient *tavily.Client
}

func NewSearchService(tavilyClient *tavily.Client) *SearchService {
	return &SearchService{
		tavilyClient: tavilyClient,
	}
}

type SearchInput struct {
	Query      string `json:"query"`
	MaxResults int    `json:"max_results,omitempty"`
	Depth      string `json:"depth,omitempty"`
}

func (s *SearchService) SearchWeb(input SearchInput) (*tavily.SearchResponse, error) {
	req := tavily.SearchRequest{
		Query:       input.Query,
		MaxResults:  input.MaxResults,
		SearchDepth: input.Depth,
	}

	return s.tavilyClient.Search(req)
}
