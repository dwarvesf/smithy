package handler

import "net/http"

var (
	ErrorNotFoundModelList = errorNotFoundModelList{}
)

type errorNotFoundModelList struct{}

func (errorNotFoundModelList) Error() string {
	return "Model list not found, please sync"
}

func (errorNotFoundModelList) StatusCode() int {
	return http.StatusBadRequest
}
