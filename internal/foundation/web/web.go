package web

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// we can enrich this functionality later let's just get this bitch up and running
type App struct {
	*chi.Mux
	// shutdown chan os.Signal
	// mw []Middleware
}

func NewApp() *App {
	r := chi.NewRouter()
	return &App{
		r,
	}
}

func (a App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Mux.ServeHTTP(w, r)
}

// add `ops []operations.Option` before middleware variadic
func (a *App) Handle(method string, path string, handler Handler, mw ...Middleware) {
	// Wrap handler specific middlwares

	// Wrap service general middlewares

	// Request execution
	h := func(w http.ResponseWriter, r *http.Request) {
		// start trace, which I'm assuming relies on upstream logic for this
		ctx := r.Context()

		// call wrapped handler
		if err := handler(ctx, w, r); err != nil {
			// signal shutdown
			return
		}
	}

	a.Mux.MethodFunc(method, path, h)
}
