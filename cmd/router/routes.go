package router

import (
	"net/http"

	"github.com/ADEMOLA200/ECommerce/cmd/models"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	App *models.Application
}

func (r *Router) Handler() http.Handler {
	mux := chi.NewRouter()

	return mux
}
