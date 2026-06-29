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

func taskDomainFromModel(taskModel TaskModel) domain.Task {
	return domain.NewTask(
		taskModel.ID,
		taskModel.Version,
		taskModel.Title,
		taskModel.Description,
		taskModel.Completed,
		taskModel.CreatedAt,
		taskModel.CompletedAt,
		taskModel.AuthorUserID,
	)
}

func taskDomainsFromModels(tasksModels []TaskModel) []domain.Task {
	domains := make([]domain.Task, len(tasksModels))

	for i, model := range tasksModels {
		domains[i] = taskDomainFromModel(model)
	}

	return domains
}
