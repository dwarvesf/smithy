package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"

	auth "github.com/dwarvesf/smithy/backend/auth"
	"github.com/dwarvesf/smithy/backend/endpoints"
)

// NewHTTPHandler http handler
func NewHTTPHandler(endpoints endpoints.Endpoints,
	logger log.Logger,
	useCORS bool, jwtSecretKey string) http.Handler {
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

	tokenAuth := jwtauth.New("HS256", []byte(jwtSecretKey), nil)

	// admin group
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Use(auth.Authorization)

		r.Post("/query", httptransport.NewServer( // Post query for case a query have more than 2048 character
			endpoints.DBQuery,
			decodeDBQueryRequest,
			httptransport.EncodeJSONResponse,
			options...,
		).ServeHTTP)

		r.Get("/agent-sync", httptransport.NewServer(
			endpoints.AgentSync,
			httptransport.NopRequestDecoder,
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

		r.Post("/hooks", httptransport.NewServer(
			endpoints.AddHook,
			decodeAddHookRequest,
			httptransport.EncodeJSONResponse,
			options...,
		).ServeHTTP)

		r.Put("/update", httptransport.NewServer(
			endpoints.DBUpdate,
			decodeDBUpdateRequest,
			httptransport.EncodeJSONResponse,
			options...,
		).ServeHTTP)

		r.Delete("/delete", httptransport.NewServer(
			endpoints.DBDelete,
			decodeDBDeleteRequest,
			httptransport.EncodeJSONResponse,
			options...,
		).ServeHTTP)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", httptransport.NewServer(
			endpoints.Login,
			decodeLoginRequest,
			httptransport.EncodeJSONResponse,
			options...,
		).ServeHTTP)
	})

	r.Get("/list-config-version", httptransport.NewServer(
		endpoints.ListVersion,
		httptransport.NopRequestDecoder,
		httptransport.EncodeJSONResponse,
		options...,
	).ServeHTTP)

	r.Post("/revert-config-version", httptransport.NewServer(
		endpoints.RevertVersion,
		decodeRevertVersion,
		httptransport.EncodeJSONResponse,
		options...,
	).ServeHTTP)

	return r
}
