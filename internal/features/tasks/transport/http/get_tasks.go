package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/gptikhomirov/go-rest-prod/internal/core/logger"
	core_http_request "github.com/gptikhomirov/go-rest-prod/internal/core/transport/http/request"
	core_http_response "github.com/gptikhomirov/go-rest-prod/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

const userIDQueryKey = "user_id"

func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := core_http_request.GetLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get limit/offset query params",
		)

		return
	}

	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryKey)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user_id query param",
		)

		return
	}

	tasks, err := h.tasksService.GetTasks(ctx, limit, offset, userID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get task list",
		)

		return
	}

	response := GetTasksResponse(taskDTOsFromDomains(tasks))

	responseHandler.JSONResponse(response, http.StatusOK)
}
