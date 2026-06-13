package developer

import "testing"

func TestAgentPlaybookNonEmpty(t *testing.T) {
	if AgentPlaybook() == "" {
		t.Fatal("expected non-empty agent playbook")
	}
	if IntegrationGuide() == "" {
		t.Fatal("expected non-empty integration guide")
	}
}
