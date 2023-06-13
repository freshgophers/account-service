package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"account-service/internal/domain/user"
	"account-service/internal/service/account"
	"account-service/pkg/server/status"
	"account-service/pkg/store"
)

type AccountHandler struct {
	accountService *account.Service
}

func NewAccountHandler(a *account.Service) *AccountHandler {
	return &AccountHandler{accountService: a}
}

func (h *AccountHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.get)
		r.Put("/", h.update)
	})

	return r
}

// Read the account from the database
//
//	@Summary	Read the account from the database
//	@Tags		accounts
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"path param"
//	@Success	200	{object}	status.Response
//	@Failure	404	{object}	status.Response
//	@Failure	500	{object}	status.Response
//	@Router		/accounts/{id} [get]
func (h *AccountHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	res, err := h.accountService.GetByID(r.Context(), id)
	if err != nil && err != store.ErrorNotFound {
		render.Render(w, r, status.InternalServerError(err))
		return
	}

	if err == store.ErrorNotFound {
		render.Render(w, r, status.NotFound(err))
		return
	}

	render.JSON(w, r, status.OK(res))
}

// Update the account in the database
//
//	@Summary	Update the account in the database
//	@Tags		accounts
//	@Accept		json
//	@Produce	json
//	@Param		id		path	string			true	"path param"
//	@Param		request	body	user.Request	true	"body param"
//	@Success	200
//	@Failure	400	{object}	status.Response
//	@Failure	404	{object}	status.Response
//	@Failure	500	{object}	status.Response
//	@Router		/accounts/{id} [put]
func (h *AccountHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	req := user.Request{}
	if err := render.Bind(r, &req); err != nil {
		render.Render(w, r, status.BadRequest(err, req))
		return
	}

	err := h.accountService.Update(r.Context(), id, req)
	if err != nil && err != store.ErrorNotFound {
		render.Render(w, r, status.InternalServerError(err))
		return
	}

	if err == store.ErrorNotFound {
		render.Render(w, r, status.NotFound(err))
		return
	}
}
