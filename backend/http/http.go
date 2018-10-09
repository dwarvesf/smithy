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

	"github.com/dwarvesf/smithy/backend/auth"
	backendConfig "github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/backend/endpoints"
	"github.com/dwarvesf/smithy/backend/service"
)

// NewHTTPHandler http handler
func NewHTTPHandler(endpoints endpoints.Endpoints,
	logger log.Logger,
	useCORS bool, cfg *backendConfig.Config, s service.Service) http.Handler {
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
		r.Use(auth.Authorization(cfg, s))

		r.Get("/agent-sync", httptransport.NewServer(
			endpoints.AgentSync,
			httptransport.NopRequestDecoder,
			httptransport.EncodeJSONResponse,
			options...,
		).ServeHTTP)

		r.Route("/databases/{db_name}", func(r chi.Router) {
			r.Route("/view", func(r chi.Router) {
				r.Post("/", httptransport.NewServer(
					endpoints.ViewAdd,
					decodeAddView,
					httptransport.EncodeJSONResponse,
					options...,
				).ServeHTTP)

				r.Get("/", httptransport.NewServer(
					endpoints.ViewList,
					decodeListView,
					httptransport.EncodeJSONResponse,
					options...,
				).ServeHTTP)

				r.Route("/{sql_id}", func(r chi.Router) {
					r.Delete("/", httptransport.NewServer(
						endpoints.ViewDelete,
						decodeDeleteView,
						httptransport.EncodeJSONResponse,
						options...,
					).ServeHTTP)

					r.Post("/execute", httptransport.NewServer(
						endpoints.ViewExecute,
						decodeExecuteView,
						httptransport.EncodeJSONResponse,
						options...,
					).ServeHTTP)
				})
			})

			r.Route("/table/{table_name}", func(r chi.Router) {
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

		r.Post("/settings/password", httptransport.NewServer(
			endpoints.ChangePassword,
			decodeChangePasswordRequest,
			httptransport.EncodeJSONResponse,
			options...,
		).ServeHTTP)

		r.Route("/groups", func(r chi.Router) {
			r.Use(auth.RequireAdmin(s))

			r.Get("/", httptransport.NewServer(
				endpoints.GroupFindAll,
				httptransport.NopRequestDecoder,
				httptransport.EncodeJSONResponse,
				options...,
			).ServeHTTP)

			r.Post("/", httptransport.NewServer(
				endpoints.GroupCreate,
				decodeCreateGroup,
				httptransport.EncodeJSONResponse,
				options...,
			).ServeHTTP)

			r.Route("/{group_id}", func(r chi.Router) {
				r.Get("/", httptransport.NewServer(
					endpoints.GroupFind,
					decodeFindGroup,
					httptransport.EncodeJSONResponse,
					options...,
				).ServeHTTP)

				r.Delete("/", httptransport.NewServer(
					endpoints.GroupDelete,
					decodeDeleteGroup,
					httptransport.EncodeJSONResponse,
					options...,
				).ServeHTTP)

				r.Put("/", httptransport.NewServer(
					endpoints.GroupUpdate,
					decodeUpdateGroup,
					httptransport.EncodeJSONResponse,
					options...,
				).ServeHTTP)

				r.Route("/permissions", func(r chi.Router) {
					r.Get("/", httptransport.NewServer(
						endpoints.PermissionFindByGroup,
						decodePermissionFindByGroup,
						httptransport.EncodeJSONResponse,
						options...,
					).ServeHTTP)
				})
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(auth.RequireAdmin(s))

			r.Get("/", httptransport.NewServer(
				endpoints.UserFindAll,
				httptransport.NopRequestDecoder,
				httptransport.EncodeJSONResponse,
				options...,
			).ServeHTTP)

			r.Post("/", httptransport.NewServer(
				endpoints.UserCreate,
				decodeCreateUser,
				httptransport.EncodeJSONResponse,
				options...,
			).ServeHTTP)

			r.Route("/{user_id}", func(r chi.Router) {
				r.Get("/", httptransport.NewServer(
					endpoints.UserFind,
					decodeFindUser,
					httptransport.EncodeJSONResponse,
					options...,
				).ServeHTTP)

				r.Delete("/", httptransport.NewServer(
					endpoints.UserDelete,
					decodeDeleteUser,
					httptransport.EncodeJSONResponse,
					options...,
				).ServeHTTP)

				r.Put("/", httptransport.NewServer(
					endpoints.UserUpdate,
					decodeUpdateUser,
					httptransport.EncodeJSONResponse,
					options...,
				).ServeHTTP)

				r.Route("/permissions", func(r chi.Router) {
					r.Get("/", httptransport.NewServer(
						endpoints.PermissionFindByUser,
						decodePermissionFindByUser,
						httptransport.EncodeJSONResponse,
						options...,
					).ServeHTTP)
				})
			})
		})

		r.Route("/permissions", func(r chi.Router) {
			r.Use(auth.RequireAdmin(s))

			r.Route("/{permission_id}", func(r chi.Router) {
				r.Put("/", httptransport.NewServer(
					endpoints.PermissionUpdate,
					decodePermissionUpdate,
					httptransport.EncodeJSONResponse,
					options...,
				).ServeHTTP)
			})
		})
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", httptransport.NewServer(
			endpoints.Login,
			decodeLoginRequest,
			httptransport.EncodeJSONResponse,
			options...,
		).ServeHTTP)

		r.Post("/identify", httptransport.NewServer(
			endpoints.FindAccount,
			decodeFindAccountRequest,
			httptransport.EncodeJSONResponse,
			options...,
		).ServeHTTP)
	})

	r.Route("/reset", func(r chi.Router) {
		r.Post("/send-mail", httptransport.NewServer(
			endpoints.SendEmail,
			decodeSendEmailRequest,
			httptransport.EncodeJSONResponse,
			options...,
		).ServeHTTP)

		r.Post("/confirm-code", httptransport.NewServer(
			endpoints.ConfirmCode,
			decodeConfirmCodeRequest,
			httptransport.EncodeJSONResponse,
			options...,
		).ServeHTTP)

		r.Post("/password", httptransport.NewServer(
			endpoints.ResetPassword,
			decodeResetPasswordRequest,
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
