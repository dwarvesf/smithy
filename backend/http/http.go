package http

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"

	auth "github.com/dwarvesf/smithy/backend/auth"
	backendConfig "github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/backend/endpoints"
)

// NewHTTPHandler http handler
func NewHTTPHandler(endpoints endpoints.Endpoints,
	logger log.Logger,
	useCORS bool, cfg *backendConfig.Config) http.Handler {
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

	tokenAuth := jwtauth.New("HS256", []byte(cfg.Authentication.SerectKey), nil)

	if os.Getenv("ENV") == "development" {
		fs := http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./swaggerui")))
		r.Get("/swaggerui", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/swaggerui/", http.StatusMovedPermanently)
		}))
		r.Get("/swaggerui/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fs.ServeHTTP(w, r)
		}))
	}

	// admin group
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(auth.Authenticator)
		r.Use(auth.Authorization(cfg))

		r.Get("/agent-sync", httptransport.NewServer(
			endpoints.AgentSync,
			httptransport.NopRequestDecoder,
			httptransport.EncodeJSONResponse,
			options...,
		).ServeHTTP)

		r.Route("/databases/{db_name}/{table_name}", func(r chi.Router) {
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

	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", httptransport.NewServer(
			endpoints.Login,
			decodeLoginRequest,
			httptransport.EncodeJSONResponse,
			options...,
		).ServeHTTP)
	})

	r.Route("/config-versions", func(r chi.Router) {
		r.Get("/", httptransport.NewServer(
			endpoints.ListVersion,
			httptransport.NopRequestDecoder,
			httptransport.EncodeJSONResponse,
			options...,
		).ServeHTTP)

		r.Post("/revert", httptransport.NewServer(
			endpoints.RevertVersion,
			decodeRevertVersion,
			httptransport.EncodeJSONResponse,
			options...,
		).ServeHTTP)
	})

	return r
}
