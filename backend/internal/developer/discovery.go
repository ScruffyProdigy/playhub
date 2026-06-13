package developer

// DiscoveryPrompt returns the agent interview script for understanding a game
// before drafting catalog metadata or seatTemplate guidance.
func DiscoveryPrompt() string {
	return `# Discover your game

Start with one open-ended prompt — do not lead with a checklist.

**Opening (ask first):**

> Tell me about the game you're thinking of — as much or as little as you have. What's the idea, how do people play together, who it's for, anything you're excited about or unsure about.

**Then clarify only what's missing:**

Read what they shared and identify gaps. You need enough to draft catalog copy and a seatTemplate plan. Ask follow-ups conversationally — one or two at a time, not a wall of questions. If they already answered something, do not re-ask it.

**Confirm what you think you know:** If you can infer an answer but aren't fully sure, check it with the developer instead of guessing or asking from scratch. For example: "It sounds like this is mostly a 2-player game — would you say that's fair?" Same for structure, vibe, session length, and the rest.

| Topic | Why it matters | Only ask if unclear |
|-------|----------------|---------------------|
| Player count | seatTemplate / game-modes | min/max, fixed or variable |
| Structure | seatTemplate | duel, free-for-all, teams, or roles |
| Social mode | tags + tone | competitive, cooperative, or party |
| Session length | tags + copy | quick rounds vs longer sessions |
| Vibe / audience | catalog voice | casual, brainy, chaotic, tactical, etc. |
| API URL | registration | public HTTPS hosting plan (not localhost) |

**Draft (show for approval, do not save yet):**
- shortDescription (~120 chars, JoinQuest tone)
- longDescription (2–4 paragraphs for the detail page)
- howToPlay (3–6 bullet steps)
- tags (1–3 IDs from catalogTagTaxonomy)
- seatTemplate guidance (point to seat-templates cookbook: duel, 3v3, composition)

Always show drafts to the developer for approval before calling updateMyGameMetadata.
`
}
