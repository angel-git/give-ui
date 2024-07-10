package main

import (
	"context"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"spt-give-ui/components"
	"spt-give-ui/spt"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

func NewChiRouter(app *App) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/connect", func(w http.ResponseWriter, r *http.Request) {
		host := r.FormValue("host")
		port := r.FormValue("port")
		serverInfo, err := spt.ConnectToSptServer(host, port)
		if err != nil {
			templ.Handler(components.ErrorConnection(err.Error())).ServeHTTP(w, r)
		}
		// store initial server info
		app.ctx = context.WithValue(app.ctx, "serverInfo", serverInfo)
		app.ctx = context.WithValue(app.ctx, "host", host)
		app.ctx = context.WithValue(app.ctx, "port", port)

		profiles, err := spt.LoadProfiles(host, port)
		if err != nil {
			templ.Handler(components.ErrorConnection(err.Error())).ServeHTTP(w, r)
		}
		templ.Handler(components.ProfileList(profiles)).ServeHTTP(w, r)
	})

	r.Get("/session/{id}", func(w http.ResponseWriter, r *http.Request) {
		sessionId := chi.URLParam(r, "id")
		app.ctx = context.WithValue(app.ctx, "sessionId", sessionId)
		allItems, err := spt.LoadItems(app.ctx.Value("host").(string), app.ctx.Value("port").(string))
		if err != nil {
			// TODO create new type of error template
			templ.Handler(components.ErrorConnection(err.Error())).ServeHTTP(w, r)
		}
		templ.Handler(components.ItemsList(allItems)).ServeHTTP(w, r)
	})

	return r
}
