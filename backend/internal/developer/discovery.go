package developer

// DiscoveryPrompt returns the agent interview script for understanding a game
// before drafting catalog metadata or seatTemplate guidance.
func DiscoveryPrompt() string {
	return `# Discover your game

Ask the developer these questions before drafting catalog copy or suggesting a seatTemplate.

1. **One-liner** — What do players do together in one sentence?
2. **Player count** — Typical group size? Min/max? Fixed teams or variable?
3. **Structure** — Head-to-head duel, free-for-all, teams, or roles (e.g. clue-giver + guessers)?
4. **Social mode** — Competitive, cooperative, or party/social?
5. **Session length** — Quick rounds (~5 min) or longer sessions?
6. **Vibe** — How should the catalog card feel (casual, brainy, chaotic, tactical)?

After answers, draft:
- shortDescription (~120 chars, JoinQuest tone)
- longDescription (2–4 paragraphs for the detail page)
- howToPlay (3–6 bullet steps)
- tags (1–3 IDs from catalogTagTaxonomy)
- seatTemplate guidance (point to seat-templates cookbook: duel, 3v3, composition)

Always show drafts to the developer for approval before calling updateMyGameMetadata.
`
}
