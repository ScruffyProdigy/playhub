package developer

import _ "embed"

//go:embed integration_guide.md
var integrationGuide string

// IntegrationGuide returns the developer integration guide markdown.
func IntegrationGuide() string {
	return integrationGuide
}
