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
	"spt-give-ui/backend/images"
	"spt-give-ui/backend/images/cache"
	"spt-give-ui/backend/images/cache_presets"
	"spt-give-ui/backend/locale"
	"spt-give-ui/backend/logger"
	"spt-give-ui/backend/models"
	"spt-give-ui/components"
	"strconv"
)

// ctx variables
const contextSessionId = "sessionId"
const contextProfiles = "profiles"
const contextAllItems = "allItems"
const contextAllBSGItems = "AllBSGItems"
const contextTraders = "traders"
const contextFavoriteSearch = "contextFavoriteSearch"
const contextServerInfo = "contextServerInfo"
const contextLocales = "contextLocales"

// App struct
type App struct {
	config       *config.Config
	ctx          context.Context
	localeMenu   *menu.Menu
	settingsMenu *menu.Menu
	menu         *menu.Menu
	name         string
	version      string
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
		app.ctx = context.WithValue(app.ctx, contextAllBSGItems, nil)
		app.ctx = context.WithValue(app.ctx, contextFavoriteSearch, false)
		app.ctx = context.WithValue(app.ctx, contextTraders, false)
		app.ctx = context.WithValue(app.ctx, contextServerInfo, nil)
		app.ctx = context.WithValue(app.ctx, contextLocales, nil)
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
		if app.ctx.Value(contextLocales) == nil {
			localeCode := locale.ConvertLocale(app.config.GetLocale())
			locales, err := api.GetLocaleFromServer(app.config.GetSptUrl(), localeCode)
			if err != nil {
				templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
				return
			}
			app.ctx = context.WithValue(app.ctx, contextLocales, locales)
		}
		locales := app.ctx.Value(contextLocales).(*models.Locales)

		var allItems *models.AllItems
		var err error

		if app.ctx.Value(contextAllItems) == nil {
			itemsResponse, err := api.LoadItems(app.config.GetSptUrl())
			if err != nil {
				templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
				return
			}
			allItems, err = api.ParseItems(itemsResponse, locales)
			if err != nil {
				templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
				return
			}
			app.ctx = context.WithValue(app.ctx, contextAllBSGItems, itemsResponse.Items)
			app.ctx = context.WithValue(app.ctx, contextAllItems, allItems)
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

		profile := getProfileFromSession(app)

		skills, err := api.LoadSkills(profile, locales)
		if err != nil {
			templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
			return
		}
		serverInfo := app.ctx.Value(contextServerInfo).(*models.ServerInfo)
		traders, err := api.LoadTraders(app.config.GetSptUrl(), profile, sessionId, locales)
		addImageToWeaponBuild(app, &profile.UserBuilds.WeaponBuilds)

		templ.Handler(components.MainPage(app.name, app.version, allItems, isFavorite, &profile, traders, skills, serverInfo.MaxLevel)).ServeHTTP(w, r)
	}
}

func getItemDetails(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		itemId := chi.URLParam(r, "id")
		allItems := app.ctx.Value(contextAllItems).(*models.AllItems)
		bsgItems := app.ctx.Value(contextAllBSGItems).(map[string]models.BSGItem)
		item := allItems.Items[itemId]
		bsgItem := bsgItems[itemId]

		globalIdx := slices.IndexFunc(allItems.GlobalPresets, func(i models.ViewPreset) bool {
			return item.Id == i.Encyclopedia
		})
		maybePresetId := ""
		if globalIdx != -1 {
			maybePresetId = allItems.GlobalPresets[globalIdx].Id
		}
		hash := cache.GetItemHash(bsgItem, bsgItems)
		imageBase64, err := loadImage(app, hash)
		if err == nil {
			item.ImageBase64 = imageBase64
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
		locales := app.ctx.Value(contextLocales).(*models.Locales)

		err := reloadProfiles(app)
		if err != nil {
			templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
			return
		}

		profile := getProfileFromSession(app)

		traders, err := api.LoadTraders(app.config.GetSptUrl(), profile, sessionId, locales)
		if err != nil {
			templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
			return
		}
		templ.Handler(components.Traders(&profile, traders)).ServeHTTP(w, r)
	}
}

func getSkills(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locales := app.ctx.Value(contextLocales).(*models.Locales)

		err := reloadProfiles(app)
		if err != nil {
			templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
			return
		}

		profile := getProfileFromSession(app)
		skills, err := api.LoadSkills(profile, locales)
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

		err := reloadProfiles(app)
		if err != nil {
			templ.Handler(getErrorComponent(app, err.Error())).ServeHTTP(w, r)
			return
		}
		profile := getProfileFromSession(app)
		weaponBuilds := profile.UserBuilds.WeaponBuilds
		addImageToWeaponBuild(app, &weaponBuilds)

		allItems := app.ctx.Value(contextAllItems).(*models.AllItems)

		templ.Handler(components.UserWeapons(allItems, weaponBuilds)).ServeHTTP(w, r)

	}
}

func getUserWeaponModal(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		weaponBuildId := chi.URLParam(r, "id")

		profile := getProfileFromSession(app)
		weaponBuilds := profile.UserBuilds.WeaponBuilds
		weaponBuildsIdx := slices.IndexFunc(weaponBuilds, func(i models.WeaponBuild) bool {
			return i.Id == weaponBuildId
		})
		weaponBuild := weaponBuilds[weaponBuildsIdx]
		addImageToWeaponBuildAttachments(app, &weaponBuild)

		allItems := app.ctx.Value(contextAllItems).(*models.AllItems)

		templ.Handler(components.UserWeaponModal(allItems, weaponBuild)).ServeHTTP(w, r)
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

func loadImage(app *App, hash int32) (string, error) {
	session := app.ctx.Value(contextSessionId).(string)
	var loader images.ImageLoader
	var url string
	cacheFolder := app.config.GetCacheFolder()
	if cacheFolder != "" {
		loader = &images.LocalImageLoader{}
		url = cacheFolder
	} else {
		loader = &images.ServerImageLoader{}
		url = app.config.GetSptUrl()
	}
	return loader.LoadImage(url, session, fmt.Sprint(hash))
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

func reloadProfiles(app *App) error {
	profiles, err := api.LoadProfiles(app.config.GetSptUrl())
	app.ctx = context.WithValue(app.ctx, contextProfiles, profiles)
	return err
}

func getProfileFromSession(app *App) models.SPTProfile {
	sessionId := app.ctx.Value(contextSessionId).(string)
	allProfiles := app.ctx.Value(contextProfiles).([]models.SPTProfile)
	profileIdx := slices.IndexFunc(allProfiles, func(i models.SPTProfile) bool {
		return i.Info.Id == sessionId
	})
	return allProfiles[profileIdx]
}

func addImageToWeaponBuild(app *App, weaponBuilds *[]models.WeaponBuild) {
	bsgItems := app.ctx.Value(contextAllBSGItems).(map[string]models.BSGItem)

	for i := range *weaponBuilds {
		weaponBuild := &(*weaponBuilds)[i]

		idx := slices.IndexFunc(*weaponBuild.Items, func(i models.WeaponBuildItem) bool {
			return i.Id == weaponBuild.Root
		})

		imageHash := cache_presets.GetItemHash((*weaponBuild.Items)[idx], *weaponBuild.Items, bsgItems)
		imageBase64, err := loadImage(app, imageHash)
		var ImageBase64 string
		if err != nil {
			ImageBase64 = ""
		} else {
			ImageBase64 = imageBase64
		}
		weaponBuild.ImageBase64 = ImageBase64
	}
}

func addImageToWeaponBuildAttachments(app *App, weaponBuild *models.WeaponBuild) {
	bsgItems := app.ctx.Value(contextAllBSGItems).(map[string]models.BSGItem)

	for j := range *weaponBuild.Items {
		weaponAttachment := &(*weaponBuild.Items)[j]
		attachmentHash := cache.GetItemHash(bsgItems[weaponAttachment.Tpl], bsgItems)
		attachmentImageBase64, err := loadImage(app, attachmentHash)
		var AttachmentImageBase64 string
		if err != nil {
			AttachmentImageBase64 = ""
		} else {
			AttachmentImageBase64 = attachmentImageBase64
		}
		weaponAttachment.ImageBase64 = AttachmentImageBase64
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
	r.Get("/user-weapons-modal/{id}", getUserWeaponModal(app))
	// this is not used as it is disabled in the template
	// https://github.com/angel-git/give-ui/issues/49
	r.Post("/magazine-loadouts/{id}", addMagazineLoadout(app))

	return r
}
