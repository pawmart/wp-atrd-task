package handler

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrResponse struct {
	HTTPStatusCode int    `json:"-"`
	Description    string `json:"description"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

var ErrInvalidInput = &ErrResponse{HTTPStatusCode: 405, Description: "Invalid input"}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, Description: "Secret not found"}

var ErrInternalServerError = &ErrResponse{HTTPStatusCode: 500, Description: "Internal server error"}
