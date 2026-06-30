package users_transport_http

import (
	"net/http"

	core_logger "github.com/gptikhomirov/go-rest-prod/internal/core/logger"
	core_http_request "github.com/gptikhomirov/go-rest-prod/internal/core/transport/http/request"
	core_http_response "github.com/gptikhomirov/go-rest-prod/internal/core/transport/http/response"
)

type GetUsersResponse []UserDTOResponse

// GetUsers godoc
// @Summary     Получить список пользователей
// @Description Получить список пользователей с пагинацией
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       limit query int false                         "Лимит записей"
// @Param       offset query int false                        "Смещение"
// @Success     200 {object} GetUsersResponse                 "Список пользователей"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /users [get]
func (h *UsersHTTPHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := core_http_request.GetLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get limit/offset query params",
		)
	}

	userDomains, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")
	}

	response := GetUsersResponse(usersDTOFromDomains(userDomains))

	responseHandler.JSONResponse(response, http.StatusOK)
}
