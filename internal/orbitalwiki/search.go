package orbitalwiki

import (
	"context"
	"errors"
)

// PART 1 — implement without AI assistance. See README.md and docs/tools/search.md.

// SearchParams holds the inputs for GET /search.
type SearchParams struct {
	// TODO(Part 1)
}

// SearchResponse is a decoded, successful GET /search response.
type SearchResponse struct {
	// TODO(Part 1)
}

// Search calls GET /search.
func (c *Client) Search(ctx context.Context, p SearchParams) (*SearchResponse, error) {
	// TODO(Part 1)
	return nil, errors.New("not implemented")
}
