package http

import (
	"account-service/internal/service/otp"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"account-service/internal/domain/secret"
	"account-service/pkg/server/status"
)

type OTPHandler struct {
	otpService *otp.Service
}

func NewOTPHandler(a *otp.Service) *OTPHandler {
	return &OTPHandler{otpService: a}
}

func (h *OTPHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.send)
	r.Post("/", h.check)

	return r
}

// Send otp code
//
//	@Summary	Send otp code
//	@Tags		otp
//	@Accept		json
//	@Produce	json
//	@Param		phone	query		string	true	"query param"
//	@Success	200		{object}	status.Response
//	@Failure	500		{object}	status.Response
//	@Router		/otp [get]
func (h *OTPHandler) send(w http.ResponseWriter, r *http.Request) {
	phone := r.URL.Query().Get("phone")
	if phone == "" {
		err := errors.New("key: cannot be blank")
		render.Render(w, r, status.BadRequest(err, nil))
		return
	}

	res, err := h.otpService.Send(r.Context(), phone)
	if err != nil {
		render.Render(w, r, status.InternalServerError(err))
		return
	}

	render.JSON(w, r, status.OK(res))
}

// Check OTP code
//
//	@Summary	Check OTP code
//	@Tags		otp
//	@Accept		json
//	@Produce	json
//	@Param		request	body		secret.Request	true	"body param"
//	@Success	200		{object}	status.Response
//	@Failure	400		{object}	status.Response
//	@Failure	500		{object}	status.Response
//	@Router		/otp [post]
func (h *OTPHandler) check(w http.ResponseWriter, r *http.Request) {
	req := secret.Request{}
	if err := render.Bind(r, &req); err != nil {
		render.Render(w, r, status.BadRequest(err, req))
		return
	}

	res, err := h.otpService.Check(r.Context(), req)
	if err != nil {
		render.Render(w, r, status.BadRequest(err, req))
		return
	}

	render.JSON(w, r, status.OK(res))
}
