package tasks_service

import (
	"context"
	"fmt"

	"github.com/gptikhomirov/go-rest-prod/internal/core/domain"
)

func (s *TasksService) GetTask(
	ctx context.Context,
	id int,
) (domain.Task, error) {
	taskDomain, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task from repository: %w", err)
	}

	return taskDomain, nil
}
