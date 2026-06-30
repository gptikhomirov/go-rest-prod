package tasks_transport_http

import (
	"time"

	"github.com/gptikhomirov/go-rest-prod/internal/core/domain"
)

type TaskDTOResponse struct {
	ID           int        `json:"id"             example:"1"`
	Version      int        `json:"version"        example:"1"`
	Title        string     `json:"title"          example:"Купить молоко"`
	Description  *string    `json:"description"     example:"2 литра, до пятницы"`
	Completed    bool       `json:"completed"      example:"false"`
	CreatedAt    time.Time  `json:"created_at"     example:"2026-06-30T12:00:00Z"`
	CompletedAt  *time.Time `json:"completed_at"    example:"2026-06-30T18:30:00Z"`
	AuthorUserID int        `json:"author_user_id" example:"1"`
}

func taskDTOFromDomain(task domain.Task) TaskDTOResponse {
	return TaskDTOResponse{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAt:    task.CreatedAt,
		CompletedAt:  task.CompletedAt,
		AuthorUserID: task.AuthorUserID,
	}
}

func taskDTOsFromDomains(tasks []domain.Task) []TaskDTOResponse {
	tasksDTO := make([]TaskDTOResponse, len(tasks))

	for i, task := range tasks {
		tasksDTO[i] = taskDTOFromDomain(task)
	}

	return tasksDTO
}
