package statistics_repository_postgres

import core_postgres_pool "github.com/gptikhomirov/go-rest-prod/internal/core/repository/postgres/pool"

type StatisticsRepository struct {
	pool core_postgres_pool.Pool
}

func NewStatisticsRepository(
	pool core_postgres_pool.Pool,
) *StatisticsRepository {
	return &StatisticsRepository{
		pool: pool,
	}
}
