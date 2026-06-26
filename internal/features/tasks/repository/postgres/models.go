package tasks_repository_postgres

import (
	"time"

	"github.com/gptikhomirov/go-rest-prod/internal/core/domain"
)

type TaskModel struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserID int
}

func taskDomainsFromModels(tasksModels []TaskModel) []domain.Task {
	domains := make([]domain.Task, len(tasksModels))

	for i, model := range tasksModels {
		domains[i] = domain.NewTask(
			model.ID,
			model.Version,
			model.Title,
			model.Description,
			model.Completed,
			model.CreatedAt,
			model.CompletedAt,
			model.AuthorUserID,
		)
	}

	return domains
}
