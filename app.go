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
	"strconv"
	"strings"
)

// ctx variables
const contextSessionId = "sessionId"
const contextProfiles = "profiles"
const contextAllItems = "allItems"
const contextTraders = "traders"
const contextFavoriteSearch = "contextFavoriteSearch"
const contextServerInfo = "contextServerInfo"

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

func getLoginPage(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		app.ctx = context.WithValue(app.ctx, contextSessionId, nil)
		app.ctx = context.WithValue(app.ctx, contextProfiles, nil)
		app.ctx = context.WithValue(app.ctx, contextAllItems, nil)
		app.ctx = context.WithValue(app.ctx, contextFavoriteSearch, false)
		app.ctx = context.WithValue(app.ctx, contextTraders, false)
		app.ctx = context.WithValue(app.ctx, contextServerInfo, nil)
		templ.Handler(components.LoginPage(app.name, app.version, app.config.GetTheme(), app.config.GetSptUrl())).ServeHTTP(w, r)
	}
}

func getProfileList(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		app.ctx = context.WithValue(app.ctx, contextServerInfo, serverInfo)

		templ.Handler(components.ProfileList(app.name, app.version, profiles)).ServeHTTP(w, r)
	}
}

func switchTheme(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.config.SwitchTheme()
	}
}

func favouriteSearch(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var isFavorite = app.ctx.Value(contextFavoriteSearch).(bool)
		app.ctx = context.WithValue(app.ctx, contextFavoriteSearch, !isFavorite)
		getMainPageForProfile(app)(w, r)
	}

}

func getMainPageForProfile(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		sessionId := chi.URLParam(r, "id")
		var isFavorite = app.ctx.Value(contextFavoriteSearch).(bool)
		app.ctx = context.WithValue(app.ctx, contextSessionId, sessionId)
		localeCode := locale.ConvertLocale(app.config.GetLocale())

		var allItems *models.AllItems
		var err error
		if app.ctx.Value(contextAllItems) == nil {
			allItems, err = api.LoadItems(app.config.GetSptUrl(), localeCode)
			if err != nil {
				templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
				return
			}
		} else {
			allItems = app.ctx.Value(contextAllItems).(*models.AllItems)
		}

		favoriteItems := app.config.GetFavoriteItems()
		for _, favoriteItem := range favoriteItems {
			item, exists := allItems.Items[favoriteItem]
			if exists {
				item.Favorite = true
				allItems.Items[favoriteItem] = item
			}
		}
		app.ctx = context.WithValue(app.ctx, contextAllItems, allItems)

		allProfiles := app.ctx.Value(contextProfiles).([]models.SPTProfile)
		allProfilesIdx := slices.IndexFunc(allProfiles, func(i models.SPTProfile) bool {
			return i.Info.Id == sessionId
		})
		profile := allProfiles[allProfilesIdx]

		skills, err := api.LoadSkills(app.config.GetSptUrl(), profile, localeCode)
		if err != nil {
			templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
			return
		}
		serverInfo := app.ctx.Value(contextServerInfo).(*models.ServerInfo)
		// trader reputation fix https://github.com/sp-tarkov/server/pull/994 is not available in versions 3.10.0, 3.10.1, 3.10.2, 3.10.3
		// TODO remove me after 3.11.0 release
		var traders []models.Trader
		if !strings.Contains(serverInfo.Version, "3.10.0") && !strings.Contains(serverInfo.Version, "3.10.1") && !strings.Contains(serverInfo.Version, "3.10.2") && !strings.Contains(serverInfo.Version, "3.10.3") {
			traders, err = api.LoadTraders(app.config.GetSptUrl(), profile, sessionId, localeCode)
			if err != nil {
				templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
				return
			}
		}
		templ.Handler(components.MainPage(app.name, app.version, allItems, isFavorite, &profile, traders, skills, serverInfo.MaxLevel)).ServeHTTP(w, r)
	}
}

func getItemDetails(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
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
	}
}

func toggleFavorite(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemId := chi.URLParam(r, "id")
		allItems := app.ctx.Value(contextAllItems).(*models.AllItems)
		app.config.ToggleFavoriteItem(itemId)

		item := allItems.Items[itemId]
		item.Favorite = !item.Favorite
		allItems.Items[itemId] = item

		app.ctx = context.WithValue(app.ctx, contextAllItems, allItems)

		getItemDetails(app)(w, r)
	}
}

func addItem(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemId := r.FormValue("id")
		amount, _ := strconv.Atoi(r.FormValue("quantity"))
		sessionId := app.ctx.Value(contextSessionId).(string)

		err := api.AddItem(app.config.GetSptUrl(), sessionId, itemId, amount)
		if err != nil {
			templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
		}
		w.Header().Set("HX-Trigger", "{\"showAddItemMessage\": \"Your item has been sent\"}")
	}
}

func updateTrader(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nickname := r.FormValue("nickname")
		spend := r.FormValue("spend")
		rep := r.FormValue("rep")
		spendOriginal := r.FormValue("spend-original")
		repOriginal := r.FormValue("rep-original")
		sessionId := app.ctx.Value(contextSessionId).(string)
		fmt.Println(rep, repOriginal)
		fmt.Println(spend, spendOriginal)

		// Convert strings to float32
		floatRep, err1 := strconv.ParseFloat(rep, 32)
		floatOriginalRep, err2 := strconv.ParseFloat(repOriginal, 32)

		if err1 != nil {
			templ.Handler(getErrorComponent(app, err1.Error())).ServeHTTP(w, r)
		}
		if err2 != nil {
			templ.Handler(getErrorComponent(app, err2.Error())).ServeHTTP(w, r)
		}

		if float32(floatRep) != float32(floatOriginalRep) {
			rep = fmt.Sprintf("%d", int(floatRep*100))

			err := api.UpdateTraderRep(app.config.GetSptUrl(), sessionId, nickname, rep)
			if err != nil {
				templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
			}
			w.Header().Set("HX-Trigger", "{\"showAddItemMessage\": \"Message sent. Don't forget to accept it\"}")
		}

		if spend != spendOriginal {
			err := api.UpdateTraderSpend(app.config.GetSptUrl(), sessionId, nickname, spend)
			if err != nil {
				templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
			}
			w.Header().Set("HX-Trigger", "{\"showAddItemMessage\": \"Message sent. Don't forget to accept it\"}")
		}
	}
}

func getTraders(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId := app.ctx.Value(contextSessionId).(string)

		profiles, err := api.LoadProfiles(app.config.GetSptUrl())
		if err != nil {
			templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
			return
		}
		app.ctx = context.WithValue(app.ctx, contextProfiles, profiles)
		profileIdx := slices.IndexFunc(profiles, func(i models.SPTProfile) bool {
			return i.Info.Id == sessionId
		})
		profile := profiles[profileIdx]
		localeCode := locale.ConvertLocale(app.config.GetLocale())

		traders, err := api.LoadTraders(app.config.GetSptUrl(), profile, sessionId, localeCode)
		if err != nil {
			templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
			return
		}
		templ.Handler(components.Traders(&profile, traders)).ServeHTTP(w, r)
	}
}

func getSkills(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId := app.ctx.Value(contextSessionId).(string)

		profiles, err := api.LoadProfiles(app.config.GetSptUrl())
		if err != nil {
			templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
			return
		}
		app.ctx = context.WithValue(app.ctx, contextProfiles, profiles)
		profileIdx := slices.IndexFunc(profiles, func(i models.SPTProfile) bool {
			return i.Info.Id == sessionId
		})
		profile := profiles[profileIdx]
		localeCode := locale.ConvertLocale(app.config.GetLocale())

		skills, err := api.LoadSkills(app.config.GetSptUrl(), profile, localeCode)
		if err != nil {
			templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
			return
		}
		serverInfo := app.ctx.Value(contextServerInfo).(*models.ServerInfo)
		templ.Handler(components.Skills(profile.Characters.PMC.InfoPMC.Level, skills, serverInfo.MaxLevel)).ServeHTTP(w, r)
	}
}
func setLevel(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId := app.ctx.Value(contextSessionId).(string)
		level, _ := strconv.Atoi(r.FormValue("level"))

		err := api.UpdateLevel(app.config.GetSptUrl(), sessionId, level)
		if err != nil {
			templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
		}
		w.Header().Set("HX-Trigger", "{\"showAddItemMessage\": \"Message sent. Don't forget to accept it\"}")

	}
}

func updateSkill(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId := app.ctx.Value(contextSessionId).(string)
		progress, _ := strconv.Atoi(r.FormValue("progress"))
		skill := r.FormValue("skill")

		err := api.UpdateSkill(app.config.GetSptUrl(), sessionId, skill, progress)
		if err != nil {
			templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
		}
		w.Header().Set("HX-Trigger", "{\"showAddItemMessage\": \"Message sent. Don't forget to accept it\"}")

	}
}

func getUserWeaponPresets(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		sessionId := app.ctx.Value(contextSessionId).(string)

		profiles, err := api.LoadProfiles(app.config.GetSptUrl())
		if err != nil {
			templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
			return
		}
		app.ctx = context.WithValue(app.ctx, contextProfiles, profiles)
		profileIdx := slices.IndexFunc(profiles, func(i models.SPTProfile) bool {
			return i.Info.Id == sessionId
		})
		weaponBuilds := profiles[profileIdx].UserBuilds.WeaponBuilds
		allItems := app.ctx.Value(contextAllItems).(*models.AllItems)

		templ.Handler(components.UserWeapons(allItems, weaponBuilds)).ServeHTTP(w, r)

	}
}

func addUserWeaponPreset(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		presetId := chi.URLParam(r, "id")
		sessionId := app.ctx.Value(contextSessionId).(string)

		err := api.AddUserWeapon(app.config.GetSptUrl(), sessionId, presetId)
		if err != nil {
			templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
		}
		w.Header().Set("HX-Trigger", "{\"showAddItemMessage\": \"Your weapon has been sent\"}")
	}
}

func addMagazineLoadout(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func NewChiRouter(app *App) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/initial", getLoginPage(app))
	r.Post("/theme", switchTheme(app))
	r.Post("/connect", getProfileList(app))
	r.Get("/connect/{id}", getMainPageForProfile(app))
	r.Get("/search/{id}", favouriteSearch(app))
	r.Get("/item/{id}", getItemDetails(app))
	r.Post("/fav/{id}", toggleFavorite(app))
	r.Post("/item", addItem(app))
	r.Post("/trader", updateTrader(app))
	r.Get("/trader", getTraders(app))
	r.Get("/skill", getSkills(app))
	r.Post("/skill", updateSkill(app))
	r.Post("/level", setLevel(app))
	r.Get("/user-weapons", getUserWeaponPresets(app))
	r.Post("/user-weapons/{id}", addUserWeaponPreset(app))
	// this is not used as it is disabled in the template
	// https://github.com/angel-git/give-ui/issues/49
	r.Post("/magazine-loadouts/{id}", addMagazineLoadout(app))

	return r
}
