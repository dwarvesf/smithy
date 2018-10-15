package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

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

func decodeAddView(ctx context.Context, r *http.Request) (interface{}, error) {
	req := endpoints.AddViewRequest{}
	dbName := chi.URLParam(r, "db_name")

	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	req.View.DatabaseName = dbName

	return req, err
}

func decodeListView(ctx context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.ListViewRequest
	dbName := chi.URLParam(r, "db_name")
	req.DatabaseName = dbName
	return req, nil
}

func decodeDeleteView(ctx context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.DeleteViewRequest
	dbName := chi.URLParam(r, "db_name")
	sqlID, err := strconv.Atoi(chi.URLParam(r, "sql_id"))
	if err != nil {
		return nil, err
	}

	req.DatabaseName = dbName
	req.SQLID = sqlID

	return req, nil
}

func decodeExecuteView(ctx context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.ExecuteViewRequest
	dbName := chi.URLParam(r, "db_name")
	sqlID, err := strconv.Atoi(chi.URLParam(r, "sql_id"))
	if err != nil {
		return nil, err
	}

	req.DatabaseName = dbName
	req.SQLID = sqlID

	return req, nil
}

func decodeCreateGroup(ctx context.Context, r *http.Request) (interface{}, error) {
	req := endpoints.CreateGroupRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	return req, err
}

func decodeFindAccountRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := endpoints.FindAccountRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	return req, err
}

func decodeDeleteGroup(ctx context.Context, r *http.Request) (interface{}, error) {
	groupID := chi.URLParam(r, "group_id")
	req := endpoints.DeleteGroupRequest{GroupID: groupID}
	return req, nil
}

func decodeReadGroup(ctx context.Context, r *http.Request) (interface{}, error) {
	groupID := chi.URLParam(r, "group_id")
	req := endpoints.ReadGroupRequest{GroupID: groupID}
	return req, nil
}

func decodeUpdateGroup(ctx context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.UpdateGroupRequest

	groupID := chi.URLParam(r, "group_id")

	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	req.GroupID = groupID

	return req, err
}

func decodeSendEmailRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := endpoints.SendEmailRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	return req, err
}

func decodeConfirmCodeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := endpoints.ConfirmCodeRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	return req, err
}

func decodeResetPasswordRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := endpoints.ResetPasswordRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	return req, err
}
