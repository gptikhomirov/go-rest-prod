package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gptikhomirov/go-rest-prod/internal/core/domain"
	core_logger "github.com/gptikhomirov/go-rest-prod/internal/core/logger"
	core_http_request "github.com/gptikhomirov/go-rest-prod/internal/core/transport/http/request"
	core_http_response "github.com/gptikhomirov/go-rest-prod/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	// Количество созданных задач за период.
	TasksCreated int `json:"tasks_created"   example:"42"`
	// Количество завершённых задач за период.
	TasksCompleted int `json:"tasks_completed" example:"30"`
	// Доля завершённых задач (0..1). null, если не было созданных задач.
	TasksCompletedRate *float64 `json:"tasks_completed_rate" example:"0.714"`
	// Среднее время выполнения задачи (Go duration). null, если нет завершённых задач.
	TasksAverageCompletionTime *string `json:"tasks_average_completion_time" example:"36h12m"`
}

// GetStatistics godoc
// @Summary     Получить статистику
// @Description Получить агрегированную статистику по задачам с опциональной фильтрацией
// @Description по пользователю и временному диапазону. Все параметры необязательны:
// @Description без них статистика считается по всем задачам за всё время.
// @Tags        statistics
// @Accept      json
// @Produce     json
// @Param       user_id query int false                       "Фильтр по ID пользователя"
// @Param       from query string false                       "Начало периода (дата, RFC3339/YYYY-MM-DD)"
// @Param       to query string false                         "Конец периода (дата, RFC3339/YYYY-MM-DD)"
// @Success     200 {object} GetStatisticsResponse            "Статистика по задачам"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /statistics [get]
func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	params, err := getUserIDFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID/from/to query params")

		return
	}

	statistics, err := h.statisticsService.GetStatistics(ctx, params.userID, params.from, params.to)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get statistics")

		return
	}

	response := toDTOFromDomain(statistics)

	responseHandler.JSONResponse(response, http.StatusOK)
}

func toDTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if statistics.TasksAverageCompletionTime != nil {
		duration := statistics.TasksAverageCompletionTime.String()
		avgTime = &duration
	}

	return GetStatisticsResponse{
		TasksCreated:               statistics.TasksCreated,
		TasksCompleted:             statistics.TasksCompleted,
		TasksCompletedRate:         statistics.TasksCompletedRate,
		TasksAverageCompletionTime: avgTime,
	}
}

type queryParams struct {
	userID *int
	from   *time.Time
	to     *time.Time
}

func getUserIDFromToQueryParams(r *http.Request) (queryParams, error) {
	const (
		userIDQueryParamKey = "user_id"
		fromQueryParamKey   = "from"
		toQueryParamKey     = "to"
	)
	emptyQueryParams := queryParams{
		userID: nil,
		from:   nil,
		to:     nil,
	}

	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return emptyQueryParams, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	from, err := core_http_request.GetDateQueryParam(r, fromQueryParamKey)
	if err != nil {
		return emptyQueryParams, fmt.Errorf("get 'from' query param: %w", err)
	}

	to, err := core_http_request.GetDateQueryParam(r, toQueryParamKey)
	if err != nil {
		return emptyQueryParams, fmt.Errorf("get 'to' query param: %w", err)
	}

	return queryParams{
		userID: userID,
		from:   from,
		to:     to,
	}, nil
}
