package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type chiRouter struct {
	mux *chi.Mux
}

//NewChiRouter construct new Router
func NewChiRouter() Router {
	return &chiRouter{
		mux: setUpChi(),
	}
}

func setUpChi() *chi.Mux {
	r := chi.NewRouter()
	setMiddlewares(r)
	return r
}

func (c *chiRouter) POST(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	c.mux.Post(uri, f)
}

func (c *chiRouter) GET(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	c.mux.Get(uri, f)
}

func setMiddlewares(r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
}
