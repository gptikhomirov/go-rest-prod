package tasks_repository_postgres

import (
	"context"
	"fmt"

	core_errors "github.com/gptikhomirov/go-rest-prod/internal/core/errors"
)

func (r *TasksRepository) DeleteTask(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM go_rest_prod.tasks
	WHERE id=$1;
	`

	commandTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf(
			"task with id=`%d`: %w",
			id,
			core_errors.ErrNotFound,
		)
	}

	return nil
}
