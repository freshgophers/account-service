package http

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"account-service/internal/domain/secret"
	"account-service/internal/service/auth"
	"account-service/pkg/server/status"
)

type AuthHandler struct {
	authService *auth.Service
}

func NewAuthHandler(a *auth.Service) *AuthHandler {
	return &AuthHandler{authService: a}
}

func (h *AuthHandler) OTPRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.sendOTP)
	r.Post("/", h.checkOTP)

	return r
}

// Send otp code
//
//	@Summary	Send otp code
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		phone	query		string	true	"query param"
//	@Success	200		{object}	status.Response
//	@Failure	500		{object}	status.Response
//	@Router		/otp [get]
func (h *AuthHandler) sendOTP(w http.ResponseWriter, r *http.Request) {
	phone := r.URL.Query().Get("phone")
	if phone == "" {
		err := errors.New("key: cannot be blank")
		render.Render(w, r, status.BadRequest(err, nil))
		return
	}

	res, err := h.authService.SendOTP(r.Context(), phone)
	if err != nil {
		render.Render(w, r, status.InternalServerError(err))
		return
	}

	render.JSON(w, r, status.OK(res))
}

// Check OTP code
//
//	@Summary	Check OTP code
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		request	body	secret.Request	true	"body param"
//	@Success	200
//	@Failure	400	{object}	status.Response
//	@Failure	500	{object}	status.Response
//	@Router		/otp [post]
func (h *AuthHandler) checkOTP(w http.ResponseWriter, r *http.Request) {
	req := secret.Request{}
	if err := render.Bind(r, &req); err != nil {
		render.Render(w, r, status.BadRequest(err, req))
		return
	}

	if err := h.authService.CheckOTP(r.Context(), req); err != nil {
		render.Render(w, r, status.BadRequest(err, req))
		return
	}

	render.JSON(w, r, status.NoContent())
}
