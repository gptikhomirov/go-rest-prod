package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	core_errors "github.com/gptikhomirov/go-rest-prod/internal/core/errors"
)

func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' not a valid int: %v: %w",
			param, key, err, core_errors.ErrInvalidArgument,
		)
	}

	return &val, nil
}

func GetDateQueryParam(r *http.Request, key string) (*time.Time, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	layout := "2006-01-02"

	date, err := time.Parse(layout, param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' not a valid date: %v: %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return &date, nil
}

func GetLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	const (
		limitKey  = "limit"
		offsetKey = "offset"
	)

	limit, err := GetIntQueryParam(r, limitKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := GetIntQueryParam(r, offsetKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return limit, offset, nil
}
