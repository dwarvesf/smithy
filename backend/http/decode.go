package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/dwarvesf/smithy/backend/endpoints"
)

func decodeDBQueryRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.DBQueryRequest
	dbName := chi.URLParam(r, "db_name")
	tableName := chi.URLParam(r, "table_name")

	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	req.SourceTable = tableName
	req.SourceDatabase = dbName

	return req, err
}

func decodeDBCreateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.DBCreateRequest
	dbName := chi.URLParam(r, "db_name")
	tableName := chi.URLParam(r, "table_name")

	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	req.TableName = tableName
	req.DatabaseName = dbName

	return req, err
}

func decodeDBUpdateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.DBUpdateRequest
	dbName := chi.URLParam(r, "db_name")
	tableName := chi.URLParam(r, "table_name")

	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	req.TableName = tableName
	req.DatabaseName = dbName

	return req, err
}

func decodeChangePasswordRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.ChangePasswordRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	return req, err
}

func decodeDBDeleteRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.DBDeleteRequest
	dbName := chi.URLParam(r, "db_name")
	tableName := chi.URLParam(r, "table_name")

	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	req.TableName = tableName
	req.DatabaseName = dbName

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
