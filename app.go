package main

import (
	"context"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"slices"
	"spt-give-ui/backend/spt"
	"spt-give-ui/components"
)

// ctx variables
const contextSessionId = "sessionId"
const contextHost = "host"
const contextPort = "port"
const contextAllItems = "allItems"
const contextServerInfo = "serverInfo"

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
		host := r.FormValue(contextHost)
		port := r.FormValue(contextPort)
		serverInfo, err := spt.ConnectToSptServer(host, port)
		if err != nil {
			templ.Handler(components.ErrorConnection(err.Error())).ServeHTTP(w, r)
			return
		}
		// store initial server info
		app.ctx = context.WithValue(app.ctx, contextServerInfo, serverInfo)
		app.ctx = context.WithValue(app.ctx, contextHost, host)
		app.ctx = context.WithValue(app.ctx, contextPort, port)

		profiles, err := spt.LoadProfiles(host, port)
		if err != nil {
			templ.Handler(components.ErrorConnection(err.Error())).ServeHTTP(w, r)
			return
		}
		templ.Handler(components.ProfileList(profiles)).ServeHTTP(w, r)
	})

	r.Get("/session/{id}", func(w http.ResponseWriter, r *http.Request) {
		sessionId := chi.URLParam(r, "id")
		app.ctx = context.WithValue(app.ctx, contextSessionId, sessionId)
		allItems, err := spt.LoadItems(app.ctx.Value(contextHost).(string), app.ctx.Value(contextPort).(string))
		if err != nil {
			// TODO create new type of error template
			templ.Handler(components.ErrorConnection(err.Error())).ServeHTTP(w, r)
			return
		}
		app.ctx = context.WithValue(app.ctx, contextAllItems, allItems)
		templ.Handler(components.ItemsList(allItems)).ServeHTTP(w, r)
	})

	r.Get("/item/{id}", func(w http.ResponseWriter, r *http.Request) {
		itemId := chi.URLParam(r, "id")
		allItems := app.ctx.Value(contextAllItems).(*spt.AllItems)
		itemIdx := slices.IndexFunc(allItems.Items, func(i spt.ViewItem) bool {
			return i.Id == itemId
		})
		item := allItems.Items[itemIdx]
		templ.Handler(components.ItemDetail(item)).ServeHTTP(w, r)

	})

	r.Post("/item/{id}", func(w http.ResponseWriter, r *http.Request) {
		itemId := chi.URLParam(r, "id")
		host := app.ctx.Value(contextHost).(string)
		port := app.ctx.Value(contextPort).(string)
		sessionId := app.ctx.Value(contextSessionId).(string)
		allItems := app.ctx.Value(contextAllItems).(*spt.AllItems)
		itemIdx := slices.IndexFunc(allItems.Items, func(i spt.ViewItem) bool {
			return i.Id == itemId
		})
		amount := allItems.Items[itemIdx].MaxStock

		spt.AddItem(host, port, sessionId, itemId, amount)
	})

	return r
}
