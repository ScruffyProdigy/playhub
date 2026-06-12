package developer

import _ "embed"

//go:embed guide.md
var integrationGuide string

// IntegrationGuide returns the developer integration guide markdown (stub for Phase A).
func IntegrationGuide() string {
	return integrationGuide
}
