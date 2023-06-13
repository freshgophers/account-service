package handler

import (
	"account-service/docs"
	"account-service/internal/config"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/swaggo/http-swagger/v2"

	_ "account-service/docs"
	"account-service/internal/handler/http"
	"account-service/internal/service/auth"
	"account-service/pkg/server/router"
)

type Dependencies struct {
	Config      config.Config
	AuthService *auth.Service
}

// Configuration is an alias for a function that will take in a pointer to a Handler and modify it
type Configuration func(h *Handler) error

// Handler is an implementation of the Handler
type Handler struct {
	dependencies Dependencies

	HTTP *chi.Mux
}

// New takes a variable amount of Configuration functions and returns a new Handler
// Each Configuration will be called in the order they are passed in
func New(d Dependencies, configs ...Configuration) (h *Handler, err error) {
	// Create the handler
	h = &Handler{
		dependencies: d,
	}

	// Apply all Configurations passed in
	for _, cfg := range configs {
		// Pass the service into the configuration function
		if err = cfg(h); err != nil {
			return
		}
	}

	return
}

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// WithHTTPHandler applies a http handler to the Handler
func WithHTTPHandler() Configuration {
	return func(h *Handler) (err error) {
		// Create the http handler, if we needed parameters, such as connection strings they could be inputted here
		h.HTTP = router.New()

		docs.SwaggerInfo.BasePath = "/api/v1"
		docs.SwaggerInfo.Host = h.dependencies.Config.HTTP.Host
		docs.SwaggerInfo.Schemes = []string{h.dependencies.Config.HTTP.Schema}

		h.HTTP.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", h.dependencies.Config.HTTP.Host)),
		))

		authHandler := http.NewAuthHandler(h.dependencies.AuthService)

		h.HTTP.Route("/api/v1", func(r chi.Router) {
			r.Mount("/otp", authHandler.OTPRoutes())
		})

		return
	}
}
