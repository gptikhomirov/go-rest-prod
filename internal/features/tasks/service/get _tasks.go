package tasks_service

import (
	"context"
	"fmt"

	"github.com/gptikhomirov/go-rest-prod/internal/core/domain"
	core_errors "github.com/gptikhomirov/go-rest-prod/internal/core/errors"
)

func (s *TasksService) GetTasks(
	ctx context.Context,
	limit *int,
	offset *int,
	userID *int,
) ([]domain.Task, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf(
			"limit must be non-negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf(
			"offset must be non-negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	tasks, err := s.tasksRepository.GetTasks(ctx, limit, offset, userID)
	if err != nil {
		return nil, fmt.Errorf("get users from repository: %w", err)
	}

	return tasks, nil
}
