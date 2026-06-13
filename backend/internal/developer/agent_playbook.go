package developer

import _ "embed"

//go:embed agent_playbook.md
var agentPlaybook string

// AgentPlaybook returns the agent integration playbook markdown.
func AgentPlaybook() string {
	return agentPlaybook
}
