package response

import (
	"encoding/json"
	"net/http"

	"tomasweigenast.com/url-shortener/models"
)

const (
	headerContentType = "content-type"
	mimeJson          = "application/json; charset=utf-8"
)

func Successful[T any](w http.ResponseWriter, data ...*T) {
	if len(data) > 1 {
		panic("cannot return more than 1 data object")
	}

	var maybeData *T
	if len(data) != 0 {
		maybeData = data[0]
	}

	res := &models.ApiResponse[T]{
		Success: true,
		Data:    maybeData,
	}

	w.Header().Set(headerContentType, mimeJson)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func Failed(w http.ResponseWriter, err models.ApiError, statusCode ...int) {
	res := &models.ApiResponse[any]{
		Success: true,
		Error:   err,
	}

	w.Header().Set(headerContentType, mimeJson)
	if len(statusCode) == 1 {
		w.WriteHeader(statusCode[0])
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(res)
}
