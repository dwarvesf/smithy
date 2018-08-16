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
