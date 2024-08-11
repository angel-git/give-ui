package main

import (
	"context"
	"fmt"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"net/http"
	"slices"
	"spt-give-ui/backend/api"
	"spt-give-ui/backend/models"
	"spt-give-ui/components"
)

// ctx variables
const contextSessionId = "sessionId"
const contextUrl = "url"
const contextProfiles = "profiles"
const contextAllItems = "allItems"
const contextServerInfo = "serverInfo"

// App struct
type App struct {
	ctx        context.Context
	language   string
	localeMenu *menu.Menu
	menu       *menu.Menu
	name       string
	version    string
}

// NewApp creates a new App application struct
func NewApp(name string, version string) *App {
	return &App{
		language: "en",
		name:     name,
		version:  version,
	}
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
		templ.Handler(components.LoginPage(app.name, app.version)).ServeHTTP(w, r)
	})

	r.Post("/connect", func(w http.ResponseWriter, r *http.Request) {
		url := r.FormValue(contextUrl)

		serverInfo, err := api.ConnectToSptServer(url)
		if err != nil {
			templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
			return
		}
		if serverInfo.ModVersion != app.version {
			templ.Handler(getErrorComponent(app, fmt.Sprintf("Wrong server mod version: %s", serverInfo.ModVersion))).ServeHTTP(w, r)
			return
		}
		// store initial server info
		app.ctx = context.WithValue(app.ctx, contextServerInfo, serverInfo)
		app.ctx = context.WithValue(app.ctx, contextUrl, url)

		profiles, err := api.LoadProfiles(url)
		if err != nil {
			templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
			return
		}
		app.ctx = context.WithValue(app.ctx, contextProfiles, profiles)

		templ.Handler(components.ProfileList(app.name, app.version, profiles)).ServeHTTP(w, r)
	})

	r.Post("/connect/{id}", func(w http.ResponseWriter, r *http.Request) {
		sessionId := chi.URLParam(r, "id")
		app.ctx = context.WithValue(app.ctx, contextSessionId, sessionId)
		locale := app.convertLocale()
		allItems, err := api.LoadItems(app.ctx.Value(contextUrl).(string), locale)
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

		templ.Handler(components.ItemsList(app.name, app.version, allItems, profile)).ServeHTTP(w, r)
	})

	r.Get("/item/{id}", func(w http.ResponseWriter, r *http.Request) {
		itemId := chi.URLParam(r, "id")
		allItems := app.ctx.Value(contextAllItems).(*models.AllItems)
		itemIdx := slices.IndexFunc(allItems.Items, func(i models.ViewItem) bool {
			return i.Id == itemId
		})
		item := allItems.Items[itemIdx]

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
		url := app.ctx.Value(contextUrl).(string)
		sessionId := app.ctx.Value(contextSessionId).(string)
		allItems := app.ctx.Value(contextAllItems).(*models.AllItems)
		itemIdx := slices.IndexFunc(allItems.Items, func(i models.ViewItem) bool {
			return i.Id == itemId
		})
		amount := allItems.Items[itemIdx].MaxStock

		err := api.AddItem(url, sessionId, itemId, amount)
		if err != nil {
			templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
		}
	})

	r.Post("/user-weapons/{id}", func(w http.ResponseWriter, r *http.Request) {
		presetId := chi.URLParam(r, "id")
		url := app.ctx.Value(contextUrl).(string)
		sessionId := app.ctx.Value(contextSessionId).(string)

		err := api.AddUserWeapon(url, sessionId, presetId)
		if err != nil {
			templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
		}
	})

	// this is not used as it is disabled in the template
	// https://github.com/angel-git/give-ui/issues/49
	r.Post("/magazine-loadouts/{id}", func(w http.ResponseWriter, r *http.Request) {
		magazineLoadoutId := chi.URLParam(r, "id")
		url := app.ctx.Value(contextUrl).(string)
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
				itemIdx := slices.IndexFunc(allItems.Items, func(i models.ViewItem) bool {
					return i.Id == item.TemplateId
				})
				amount := allItems.Items[itemIdx].MaxStock
				err := api.AddItem(url, sessionId, item.TemplateId, amount)
				if err != nil {
					templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
					break
				}
			}
		}
	})

	return r
}

func (a *App) setLocale(data *menu.CallbackData) {
	if a.language == data.MenuItem.Label {
		return
	}
	a.language = data.MenuItem.Label
	for _, localeMenu := range a.localeMenu.Items {
		localeMenu.Checked = false
	}
	data.MenuItem.Checked = true

	// refresh menu with the selected locale
	runtime.MenuSetApplicationMenu(a.ctx, a.menu)
	runtime.MenuUpdateApplicationMenu(a.ctx)

	// refresh to main screen
	runtime.WindowReloadApp(a.ctx)
}

func (a *App) convertLocale() string {
	switch a.language {
	case "English":
		return "en"
	case "Czech":
		return "cz"
	case "French":
		return "fr"
	case "German":
		return "ge"
	case "Hungarian":
		return "hu"
	case "Italian":
		return "it"
	case "Japanese":
		return "jp"
	case "Korean":
		return "kr"
	case "Polish":
		return "pl"
	case "Portuguese":
		return "po"
	case "Slovak":
		return "sk"
	case "Spanish":
		return "es"
	case "Spanish - Mexico":
		return "es-mx"
	case "Turkish":
		return "tu"
	case "Romanian":
		return "ro"
	case "Русский":
		return "ru"
	default:
		return "en"
	}
}
