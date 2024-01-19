package senzingrestservice

// ----------------------------------------------------------------------------
// Structs
// ----------------------------------------------------------------------------

type GenericSenzingResult struct {
	Meta  any `json:"meta"`
	Links any `json:"links"`
	Data  any `json:"data"`
}

type SearchResults struct {
	SearchResults []SearchResult `json:"searchResults"`
}

type SearchResult struct {
	EntityId   int `json:"entityId"`
	EntityName int `json:"entityId"`
}
