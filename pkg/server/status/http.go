package status

import (
	"github.com/go-chi/render"
	"net/http"
)

type Response struct {
	Status  int    `json:"-"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func (e *Response) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(e.Status)
	return nil
}

func OK(data any) render.Renderer {
	return &Response{
		Status:  http.StatusOK,
		Success: true,
		Data:    data,
	}
}

func NoContent() render.Renderer {
	return &Response{
		Status:  http.StatusOK,
		Success: true,
	}
}

func BadRequest(err error, data any) render.Renderer {
	return &Response{
		Status:  http.StatusBadRequest,
		Success: false,
		Message: err.Error(),
		Data:    data,
	}
}

func NotFound(err error) render.Renderer {
	return &Response{
		Status:  http.StatusNotFound,
		Success: false,
		Message: err.Error(),
	}
}

func InternalServerError(err error) render.Renderer {
	return &Response{
		Status:  http.StatusInternalServerError,
		Success: false,
		Message: err.Error(),
	}
}
