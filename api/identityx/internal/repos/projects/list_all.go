package projects

import (
	"context"

	"IdentityX/models"

	"lib/database"
)

// ListAll returns every project — the Key-lifecycle module crosses it so
// EnsureAll can provision every scope (platform + each project).
func (repo *Repo) ListAll(ctx context.Context) ([]models.Project, error) {
	rows, err := database.Queries(ctx, repo.q).ListProjects(ctx)
	if err != nil {
		return nil, repo.dbe(err)
	}
	out := make([]models.Project, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapProject(row))
	}
	return out, nil
}
