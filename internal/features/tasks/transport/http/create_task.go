package tasks_transport_http

import (
	"net/http"

	"github.com/gptikhomirov/go-rest-prod/internal/core/domain"
	core_logger "github.com/gptikhomirov/go-rest-prod/internal/core/logger"
	core_http_request "github.com/gptikhomirov/go-rest-prod/internal/core/transport/http/request"
	core_http_response "github.com/gptikhomirov/go-rest-prod/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	// Заголовок задачи. Обязателен, от 1 до 100 символов.
	Title string `json:"title"       validate:"required,min=1,max=100"   example:"Купить молоко"`
	// Описание задачи. Опционально; если передано — от 3 до 1000 символов.
	Description *string `json:"description" validate:"omitempty,min=3,max=1000" example:"2 литра, до пятницы"`
	// ID пользователя-автора задачи. Обязателен.
	UserID int `json:"user_id"     validate:"required"                 example:"1"`
}

type CreateTaskResponse TaskDTOResponse

// CreateTask godoc
// @Summary     Создать задачу
// @Description Создать новую задачу в системе
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       request body CreateTaskRequest true           "CreateTask тело запроса"
// @Success     201 {object} CreateTaskResponse               "Успешно созданная задача"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     404 {object} core_http_response.ErrorResponse "User not found"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /tasks [post]
func (h *TasksHTTPHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request CreateTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")

		return
	}

	taskDomain := domain.NewTaskUninitialized(
		request.Title,
		request.Description,
		request.UserID,
	)
	taskDomain, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create task")

		return
	}

	response := CreateTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}
