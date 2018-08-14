package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/dwarvesf/smithy/backend/endpoints"
)

// NewHTTPHandler http handler
func NewHTTPHandler(endpoints endpoints.Endpoints,
	logger log.Logger,
	useCORS bool) http.Handler {
	r := chi.NewRouter()

	if useCORS {
		cors := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
		})
		r.Use(cors.Handler)
	}

	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Get("/_warm", httptransport.NewServer(
		endpoint.Nop,
		httptransport.NopRequestDecoder,
		httptransport.EncodeJSONResponse,
		options...,
	).ServeHTTP)

	r.Get("/agent-sync", httptransport.NewServer(
		endpoints.AgentSync,
		httptransport.NopRequestDecoder,
		httptransport.EncodeJSONResponse,
		options...,
	).ServeHTTP)

	r.Get("/query", httptransport.NewServer(
		endpoints.DBQuery,
		decodeDBQueryRequest,
		httptransport.EncodeJSONResponse,
		options...,
	).ServeHTTP)

	r.Post("/query", httptransport.NewServer( // Post query for case a query have more than 2048 character
		endpoints.DBQuery,
		decodeDBQueryRequest,
		httptransport.EncodeJSONResponse,
		options...,
	).ServeHTTP)

	r.Post("/create", httptransport.NewServer(
		endpoints.DBCreate,
		decodeDBCreateRequest,
		httptransport.EncodeJSONResponse,
		options...,
	).ServeHTTP)

	r.Get("/models", httptransport.NewServer(
		endpoints.AvailableModels,
		httptransport.NopRequestDecoder,
		httptransport.EncodeJSONResponse,
		options...,
	).ServeHTTP)

	r.Post("/update", httptransport.NewServer(
		endpoints.DBUpdate,
		decodeDBUpdateRequest,
		httptransport.EncodeJSONResponse,
		options...,
	).ServeHTTP)

	return r
}
