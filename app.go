package main

import (
	"context"
	"fmt"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"net/http"
	"slices"
	"spt-give-ui/backend/api"
	"spt-give-ui/backend/config"
	"spt-give-ui/backend/locale"
	"spt-give-ui/backend/logger"
	"spt-give-ui/backend/models"
	"spt-give-ui/components"
)

// ctx variables
const contextSessionId = "sessionId"
const contextProfiles = "profiles"
const contextAllItems = "allItems"

// App struct
type App struct {
	config     *config.Config
	ctx        context.Context
	localeMenu *menu.Menu
	menu       *menu.Menu
	name       string
	version    string
}

// NewApp creates a new App application struct
func NewApp(name string, version string) *App {
	logger.SetupLogger()
	a := &App{
		name:    name,
		version: version,
	}
	a.config = config.LoadConfig()
	return a
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
	err := a.config.Close()
	if err != nil {
		panic("Can't close connection to database: " + err.Error())
	}
}

func getErrorComponent(app *App, err string) templ.Component {
	giveUiError := models.GiveUiError{
		AppName:    app.name,
		AppVersion: app.version,
		Error:      err,
	}
	return components.ErrorConnection(giveUiError)
}

func NewChiRouter(app *App) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/initial", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		templ.Handler(components.LoginPage(app.name, app.version, app.config.GetTheme(), app.config.GetSptUrl())).ServeHTTP(w, r)
	})

	r.Post("/theme", func(w http.ResponseWriter, r *http.Request) {
		app.config.SwitchTheme()
	})

	r.Post("/connect", func(w http.ResponseWriter, r *http.Request) {
		url := r.FormValue("url")
		app.config.SetSptUrl(url)
		serverInfo, err := api.ConnectToSptServer(url)
		if err != nil {
			templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
			return
		}
		if serverInfo.ModVersion != app.version {
			templ.Handler(getErrorComponent(app, fmt.Sprintf("Wrong server mod version: %s", serverInfo.ModVersion))).ServeHTTP(w, r)
			return
		}

		profiles, err := api.LoadProfiles(url)
		if err != nil {
			templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
			return
		}
		app.ctx = context.WithValue(app.ctx, contextProfiles, profiles)

		templ.Handler(components.ProfileList(app.name, app.version, profiles)).ServeHTTP(w, r)
	})

	r.Get("/connect/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		sessionId := chi.URLParam(r, "id")
		app.ctx = context.WithValue(app.ctx, contextSessionId, sessionId)
		localeCode := locale.ConvertLocale(app.config.GetLocale())
		allItems, err := api.LoadItems(app.config.GetSptUrl(), localeCode)
		if err != nil {
			templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
			return
		}
		app.ctx = context.WithValue(app.ctx, contextAllItems, allItems)

		allProfiles := app.ctx.Value(contextProfiles).([]models.SPTProfile)
		allProfilesIdx := slices.IndexFunc(allProfiles, func(i models.SPTProfile) bool {
			return i.Info.Id == sessionId
		})
		profile := allProfiles[allProfilesIdx]

		if r.Header.Get("HX-Request") == "true" {
			templ.Handler(components.MainPage(app.name, app.version, allItems, &profile)).ServeHTTP(w, r)
		} else {
			templ.Handler(components.MainPageWrapped(app.name, app.version, app.config.GetTheme(), allItems, &profile)).ServeHTTP(w, r)
		}

	})

	r.Get("/item/{id}", func(w http.ResponseWriter, r *http.Request) {
		itemId := chi.URLParam(r, "id")
		allItems := app.ctx.Value(contextAllItems).(*models.AllItems)
		item := allItems.Items[itemId]

		globalIdx := slices.IndexFunc(allItems.GlobalPresets, func(i models.ViewPreset) bool {
			return item.Id == i.Encyclopedia
		})
		maybePresetId := ""
		if globalIdx != -1 {
			maybePresetId = allItems.GlobalPresets[globalIdx].Id
		}

		templ.Handler(components.ItemDetail(item, maybePresetId)).ServeHTTP(w, r)

	})

	r.Post("/item/{id}", func(w http.ResponseWriter, r *http.Request) {
		itemId := chi.URLParam(r, "id")
		sessionId := app.ctx.Value(contextSessionId).(string)
		allItems := app.ctx.Value(contextAllItems).(*models.AllItems)
		amount := allItems.Items[itemId].MaxStock

		err := api.AddItem(app.config.GetSptUrl(), sessionId, itemId, amount)
		if err != nil {
			templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
		}
	})

	r.Post("/user-weapons/{id}", func(w http.ResponseWriter, r *http.Request) {
		presetId := chi.URLParam(r, "id")
		sessionId := app.ctx.Value(contextSessionId).(string)

		err := api.AddUserWeapon(app.config.GetSptUrl(), sessionId, presetId)
		if err != nil {
			templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
		}
	})

	// this is not used as it is disabled in the template
	// https://github.com/angel-git/give-ui/issues/49
	r.Post("/magazine-loadouts/{id}", func(w http.ResponseWriter, r *http.Request) {
		magazineLoadoutId := chi.URLParam(r, "id")
		sessionId := app.ctx.Value(contextSessionId).(string)
		allItems := app.ctx.Value(contextAllItems).(*models.AllItems)

		allProfiles := app.ctx.Value(contextProfiles).([]models.SPTProfile)
		allProfilesIdx := slices.IndexFunc(allProfiles, func(i models.SPTProfile) bool {
			return i.Info.Id == sessionId
		})

		magazineBuilds := allProfiles[allProfilesIdx].UserBuilds.MagazineBuilds
		magazineBuildsIdx := slices.IndexFunc(magazineBuilds, func(i models.MagazineBuild) bool {
			return i.Id == magazineLoadoutId
		})

		for _, item := range magazineBuilds[magazineBuildsIdx].Items {
			if item.TemplateId != "" {
				amount := allItems.Items[item.TemplateId].MaxStock
				err := api.AddItem(app.config.GetSptUrl(), sessionId, item.TemplateId, amount)
				if err != nil {
					templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
					break
				}
			}
		}
	})

	return r
}
