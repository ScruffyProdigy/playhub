package developer

import "fmt"

// CatalogTag describes one allowed catalog tag ID.
type CatalogTag struct {
	ID          string
	Label       string
	Description string
}

// CatalogTagTaxonomy is the fixed v1 tag vocabulary for catalog cards.
var CatalogTagTaxonomy = []CatalogTag{
	{ID: "competitive", Label: "Competitive", Description: "Winners and losers, direct opposition"},
	{ID: "cooperative", Label: "Co-op", Description: "Players win or lose together"},
	{ID: "party", Label: "Party", Description: "Social groups, casual fun"},
	{ID: "1v1", Label: "1v1", Description: "Exactly two players head-to-head"},
	{ID: "quick", Label: "Quick", Description: "Sessions under about ten minutes"},
	{ID: "words", Label: "Words", Description: "Word and language puzzles"},
	{ID: "strategy", Label: "Strategy", Description: "Planning, hidden info, tactics"},
	{ID: "casual", Label: "Casual", Description: "Low pressure, easy to pick up"},
}

// ValidateCatalogTags ensures every tag is in the taxonomy (empty slice allowed).
func ValidateCatalogTags(tags []string) error {
	if len(tags) == 0 {
		return nil
	}
	allowed := make(map[string]struct{}, len(CatalogTagTaxonomy))
	for _, t := range CatalogTagTaxonomy {
		allowed[t.ID] = struct{}{}
	}
	for _, tag := range tags {
		if _, ok := allowed[tag]; !ok {
			return fmt.Errorf("developer: unknown catalog tag %q", tag)
		}
	}
	return nil
}
