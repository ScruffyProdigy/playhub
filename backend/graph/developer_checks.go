package graph

import (
	"context"

	"github.com/scruffyprodigy/playhub/graph/model"
	"github.com/scruffyprodigy/playhub/internal/integrationchecks"
	"github.com/scruffyprodigy/playhub/internal/store"
)

func (r *Resolver) checkRunner() *integrationchecks.Runner {
	return integrationchecks.NewRunner(r.manifestFetcher())
}

func (r *Resolver) persistGameChecks(ctx context.Context, game *store.Game) ([]*model.GameIntegrationCheck, error) {
	st, err := r.requireStore()
	if err != nil {
		return nil, err
	}
	runner := r.checkRunner()
	all := append(runner.RunManifestChecks(ctx, game), integrationchecks.StubChecks()...)
	out := make([]*model.GameIntegrationCheck, 0, len(all))
	for _, check := range all {
		var message *string
		if check.Message != "" {
			msg := check.Message
			message = &msg
		}
		stored, err := st.UpsertIntegrationCheck(ctx, game.ID, check.CheckID, check.Status, message, check.Detail)
		if err != nil {
			return nil, err
		}
		out = append(out, ToGraphQLIntegrationCheck(stored))
	}
	return out, nil
}
