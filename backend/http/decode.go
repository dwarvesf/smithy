package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/schema"

	"github.com/dwarvesf/smithy/backend/endpoints"
)

func decodeDBQueryRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	decoder := schema.NewDecoder()
	req := endpoints.DBQueryRequest{}
	if err := decoder.Decode(&req, r.Form); err != nil {
		return nil, err
	}

	if err := req.UpdateColumnsByCols(); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeDBCreateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.DBCreateRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	return req, err
}

func decodeDBUpdateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.DBUpdateRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	return req, err
}

func decodeDBDeleteRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.DBDeleteRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	return req, err
}

func decodeRevertVersion(ctx context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.RevertVersionResquest

	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	return req, err
}

func decodeLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := endpoints.LoginRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	return req, err
}

func decodeAddHookRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := endpoints.AddHookRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	return req, err
}
