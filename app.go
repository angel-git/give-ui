package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"spt-give-ui/backend/api"
	"spt-give-ui/backend/config"
	client "spt-give-ui/backend/http"
	"spt-give-ui/backend/images"
	"spt-give-ui/backend/images/cache"
	"spt-give-ui/backend/images/cache_presets"
	"spt-give-ui/backend/locale"
	"spt-give-ui/backend/models"
	"spt-give-ui/backend/util"
	"spt-give-ui/components"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ctx variables
const contextSessionId = "sessionId"
const contextProfiles = "profiles"
const contextAllItems = "allItems"
const contextAllBSGItems = "AllBSGItems"
const contextFavoriteSearch = "contextFavoriteSearch"
const contextServerInfo = "contextServerInfo"
const contextLocales = "contextLocales"
const contextAllQuests = "contextAllQuests"

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
	a := &App{
		name:    name,
		version: version,
	}
	a.config = config.LoadConfig()
	client.NewClient(a.config.GetTimeoutSeconds())
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

func redirectToErrorPage(app *App, err string) {
	giveUiError := models.GiveUiError{
		AppName:    app.name,
		AppVersion: app.version,
		Error:      err,
	}
	a := components.ErrorConnection(giveUiError)
	var sb strings.Builder
	a.Render(app.ctx, &sb)
	runtime.EventsEmit(app.ctx, "error", sb.String())

}

func getLoginPage(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		app.ctx = context.WithValue(app.ctx, contextSessionId, nil)
		app.ctx = context.WithValue(app.ctx, contextProfiles, nil)
		app.ctx = context.WithValue(app.ctx, contextAllItems, nil)
		app.ctx = context.WithValue(app.ctx, contextAllBSGItems, nil)
		app.ctx = context.WithValue(app.ctx, contextFavoriteSearch, false)
		app.ctx = context.WithValue(app.ctx, contextServerInfo, nil)
		app.ctx = context.WithValue(app.ctx, contextLocales, nil)
		app.ctx = context.WithValue(app.ctx, contextAllQuests, nil)
		templ.Handler(components.LoginPage(app.name, app.version, app.config.GetTheme(), app.config.GetSptUrl())).ServeHTTP(w, r)
	}
}

func getProfileList(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		sptUrl := r.FormValue("url")
		app.config.SetSptUrl(sptUrl)
		serverInfo, err := api.ConnectToSptServer(sptUrl)
		if err != nil {
			redirectToErrorPage(app, err.Error())
			return
		}
		if serverInfo.ModVersion != app.version {
			redirectToErrorPage(app, fmt.Sprintf("Wrong server mod version: %s", serverInfo.ModVersion))
			return
		}

		profiles, err := api.LoadProfiles(sptUrl)
		if err != nil {
			redirectToErrorPage(app, err.Error())
			return
		}
		app.ctx = context.WithValue(app.ctx, contextProfiles, profiles)
		app.ctx = context.WithValue(app.ctx, contextServerInfo, serverInfo)

		templ.Handler(components.ProfileList(app.name, app.version, profiles)).ServeHTTP(w, r)
	}
}

func goToProfileList(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		serverInfo, err := api.ConnectToSptServer(app.config.GetSptUrl())
		if err != nil {
			redirectToErrorPage(app, err.Error())
			return
		}
		if serverInfo.ModVersion != app.version {
			redirectToErrorPage(app, fmt.Sprintf("Wrong server mod version: %s", serverInfo.ModVersion))
			return
		}

		profiles, err := api.LoadProfiles(app.config.GetSptUrl())
		if err != nil {
			redirectToErrorPage(app, err.Error())
			return
		}
		app.ctx = context.WithValue(app.ctx, contextProfiles, profiles)
		app.ctx = context.WithValue(app.ctx, contextServerInfo, serverInfo)

		runtime.EventsEmit(app.ctx, "clean_profile")

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
				redirectToErrorPage(app, err.Error())
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
				redirectToErrorPage(app, err.Error())
				return
			}
			allItems, err = api.ParseItems(itemsResponse, locales)
			if err != nil {
				redirectToErrorPage(app, err.Error())
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

		if app.ctx.Value(contextAllQuests) == nil {
			bsgQuests, e := api.LoadQuests(app.config.GetSptUrl(), sessionId)
			if e != nil {
				redirectToErrorPage(app, e.Error())
				return
			}
			app.ctx = context.WithValue(app.ctx, contextAllQuests, bsgQuests)
		}

		profile := getProfileFromSession(app)

		skills, err := api.LoadSkills(profile, locales)
		if err != nil {
			redirectToErrorPage(app, err.Error())
			return
		}
		serverInfo := app.ctx.Value(contextServerInfo).(*models.ServerInfo)
		traders, err := api.LoadTraders(app.config.GetSptUrl(), profile, sessionId, locales)
		addImageToWeaponBuild(app, &profile.UserBuilds.WeaponBuilds)
		addUIPropertiesToInventoryItems(app, profile.Characters.PMC.Inventory.Stash, &profile.Characters.PMC.Inventory.Items)
		quests := getCurrentActiveQuests(app, profile)
		templ.Handler(components.MainPage(app.name, app.version, allItems, isFavorite, &profile, traders, skills, serverInfo, quests)).ServeHTTP(w, r)
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
		var hash int32
		if maybePresetId != "" {
			hash = cache_presets.GetItemHash(allItems.GlobalPresets[globalIdx].Items[0], allItems.GlobalPresets[globalIdx].Items, bsgItems)
		} else {
			hash = cache.GetItemHash(bsgItem, bsgItems)
		}
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
			redirectToErrorPage(app, err.Error())
		}
		runtime.EventsEmit(app.ctx, "toast.info", "Your item has been sent")
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

		if rep != repOriginal {
			floatRep, err := strconv.ParseFloat(rep, 64)
			if err != nil {
				redirectToErrorPage(app, err.Error())
				return
			}
			rep = fmt.Sprintf("%d", int(floatRep*100))
			err = api.UpdateTraderRep(app.config.GetSptUrl(), sessionId, nickname, rep)
			if err != nil {
				redirectToErrorPage(app, err.Error())
				return
			}
			runtime.EventsEmit(app.ctx, "toast.info", "Message sent. Don't forget to accept it")
		}

		if spend != spendOriginal {
			err := api.UpdateTraderSpend(app.config.GetSptUrl(), sessionId, nickname, spend)
			if err != nil {
				redirectToErrorPage(app, err.Error())
				return
			}
			runtime.EventsEmit(app.ctx, "toast.info", "Message sent. Don't forget to accept it")
		}
	}
}

func getTraders(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		sessionId := app.ctx.Value(contextSessionId).(string)
		locales := app.ctx.Value(contextLocales).(*models.Locales)

		err := reloadProfiles(app)
		if err != nil {
			redirectToErrorPage(app, err.Error())
			return
		}

		profile := getProfileFromSession(app)

		traders, err := api.LoadTraders(app.config.GetSptUrl(), profile, sessionId, locales)
		if err != nil {
			redirectToErrorPage(app, err.Error())
			return
		}
		templ.Handler(components.Traders(&profile, traders)).ServeHTTP(w, r)
	}
}

func getSkills(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		locales := app.ctx.Value(contextLocales).(*models.Locales)

		err := reloadProfiles(app)
		if err != nil {
			redirectToErrorPage(app, err.Error())
			return
		}

		profile := getProfileFromSession(app)
		skills, err := api.LoadSkills(profile, locales)
		if err != nil {
			redirectToErrorPage(app, err.Error())
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
			redirectToErrorPage(app, err.Error())
			return
		}
		runtime.EventsEmit(app.ctx, "toast.info", "Message sent. Don't forget to accept it")

	}
}

func getFile(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId := app.ctx.Value(contextSessionId).(string)
		imageUrl := r.URL.Query().Get("url")
		imageUrlUnescape, err := url.QueryUnescape(imageUrl)
		if err != nil {
			redirectToErrorPage(app, err.Error())
			return
		}
		image, err := api.LoadFile(app.config.GetSptUrl(), sessionId, imageUrlUnescape)
		if err != nil {
			runtime.LogWarning(app.ctx, "Couldn't find avatar image for: "+imageUrlUnescape)
			return
		}
		w.Write(image)
	}
}

func getLinkedSearchModal(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemId := chi.URLParam(r, "id")
		allItems := app.ctx.Value(contextAllItems).(*models.AllItems)
		bsgItems := app.ctx.Value(contextAllBSGItems).(map[string]models.BSGItem)

		baseItem, found := bsgItems[itemId]
		if !found {
			runtime.LogWarning(app.ctx, "Couldn't find item in BSG items: "+itemId)
			runtime.EventsEmit(app.ctx, "toast.error", "Couldn't find item in BSG items")
			return
		}
		var linkedItems []models.ViewItem
		linkedItemIds := util.SearchLink(baseItem)
		for _, id := range linkedItemIds {
			item, exists := allItems.Items[id]
			if exists {
				hash := cache.GetItemHash(baseItem, bsgItems)
				imageBase64, err := loadImage(app, hash)
				if err == nil {
					item.ImageBase64 = imageBase64
				}
				linkedItems = append(linkedItems, item)
			}
		}
		templ.Handler(components.LinkedSearchModal(linkedItems)).ServeHTTP(w, r)
	}
}

func sendSptMessage(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId := app.ctx.Value(contextSessionId).(string)
		message := r.FormValue("message")

		var err error
		switch message {
		case "summer":
			err = api.SetSummerSeason(app.config.GetSptUrl(), sessionId)
		case "halloween":
			err = api.SetHalloweenSeason(app.config.GetSptUrl(), sessionId)
		case "winter":
			err = api.SetWinterSeason(app.config.GetSptUrl(), sessionId)
		case "christmas":
			err = api.SetChristmasSeason(app.config.GetSptUrl(), sessionId)
		case "stash":
			err = api.AddRowsToStash(app.config.GetSptUrl(), sessionId)
		default:
			err = api.SendGift(app.config.GetSptUrl(), sessionId, message)
		}

		if err != nil {
			redirectToErrorPage(app, err.Error())
			return
		}
		runtime.EventsEmit(app.ctx, "toast.info", "Message sent. Read the response in Tarkov dialogues")
	}
}

func updateSkill(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId := app.ctx.Value(contextSessionId).(string)
		progress, _ := strconv.Atoi(r.FormValue("progress"))
		skill := r.FormValue("skill")

		err := api.UpdateSkill(app.config.GetSptUrl(), sessionId, skill, progress)
		if err != nil {
			redirectToErrorPage(app, err.Error())
			return
		}
		runtime.EventsEmit(app.ctx, "toast.info", "Message sent. Don't forget to accept it")
	}
}

func getUserWeaponPresets(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")

		err := reloadProfiles(app)
		if err != nil {
			redirectToErrorPage(app, err.Error())
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
			redirectToErrorPage(app, err.Error())
			return
		}
		runtime.EventsEmit(app.ctx, "toast.info", "Your weapon has been sent")
	}
}

func loadImage(app *App, hash int32) (string, error) {
	session := app.ctx.Value(contextSessionId).(string)
	var loader images.ImageLoader
	var url string
	cacheFolder := app.config.GetCacheFolder()
	if !app.config.GetUseCache() {
		return "", errors.New("Cache is disabled")
	}
	if cacheFolder != "" {
		loader = &images.LocalImageLoader{}
		url = cacheFolder
	} else {
		loader = &images.ServerImageLoader{}
		url = app.config.GetSptUrl()
	}
	return loader.LoadImage(url, session, fmt.Sprint(hash))
}

func getKit(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gearId := chi.URLParam(r, "id")
		sessionId := app.ctx.Value(contextSessionId).(string)

		allProfiles := app.ctx.Value(contextProfiles).([]models.SPTProfile)
		allProfilesIdx := slices.IndexFunc(allProfiles, func(i models.SPTProfile) bool {
			return i.Info.Id == sessionId
		})

		equipmentBuilds := allProfiles[allProfilesIdx].UserBuilds.EquipmentBuilds
		equipmentBuildsIdx := slices.IndexFunc(equipmentBuilds, func(i models.EquipmentBuild) bool {
			return i.Id == gearId
		})
		equipmentBuild := equipmentBuilds[equipmentBuildsIdx]
		var slotsWithImages = []string{"Earpiece", "Headwear", "FaceCover", "ArmBand", "ArmorVest", "Eyewear", "FirstPrimaryWeapon", "Holster", "SecondPrimaryWeapon", "Scabbard", "TacticalVest", "Backpack", "SecuredContainer"}
		for _, slotWithImage := range slotsWithImages {
			addImageToKit(app, slotWithImage, equipmentBuild)
		}

		templ.Handler(components.Kit(equipmentBuilds[equipmentBuildsIdx])).ServeHTTP(w, r)
	}
}

func getKits(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		sessionId := app.ctx.Value(contextSessionId).(string)

		err := reloadProfiles(app)
		if err != nil {
			redirectToErrorPage(app, err.Error())
			return
		}
		profile := getProfileFromSession(app)
		equipmentBuilds := profile.UserBuilds.EquipmentBuilds

		templ.Handler(components.Kits(equipmentBuilds, sessionId)).ServeHTTP(w, r)
	}
}

func addImageToKit(app *App, slot string, equipmentBuild models.EquipmentBuild) {
	index := slices.IndexFunc(equipmentBuild.Items, func(i models.ItemWithUpd) bool {
		return i.SlotID != nil && *i.SlotID == slot
	})
	if index != -1 {
		image64 := calculateImageBase64FromItems(app, equipmentBuild.Items, index, app.ctx.Value(contextAllBSGItems).(map[string]models.BSGItem))
		equipmentBuild.Items[index].ImageBase64 = image64
	}
}

func addKit(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		presetId := chi.URLParam(r, "presetId")
		sessionId := app.ctx.Value(contextSessionId).(string)

		err := api.AddGearPreset(app.config.GetSptUrl(), sessionId, presetId)
		if err != nil {
			redirectToErrorPage(app, err.Error())
			return
		}
		runtime.EventsEmit(app.ctx, "toast.info", "Your kit has been sent")
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
					redirectToErrorPage(app, err.Error())
					break
				}
			}
		}
	}
}

func getStash(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		err := reloadProfiles(app)
		if err != nil {
			redirectToErrorPage(app, err.Error())
			return
		}
		profile := getProfileFromSession(app)
		addUIPropertiesToInventoryItems(app, profile.Characters.PMC.Inventory.Stash, &profile.Characters.PMC.Inventory.Items)
		templ.Handler(components.Stash(&profile)).ServeHTTP(w, r)
	}
}

func addStashItem(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemId := r.FormValue("id")
		sessionId := app.ctx.Value(contextSessionId).(string)

		err := api.AddStashItem(app.config.GetSptUrl(), sessionId, itemId)
		if err != nil {
			redirectToErrorPage(app, err.Error())
			return
		}
		runtime.EventsEmit(app.ctx, "toast.info", "Your item has been sent")
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

		idx := slices.IndexFunc(*weaponBuild.Items, func(i models.ItemWithUpd) bool {
			return i.Id == weaponBuild.Root
		})
		var ImageBase64 = calculateImageBase64FromItems(app, *weaponBuild.Items, idx, bsgItems)
		weaponBuild.ImageBase64 = ImageBase64
	}
}

func calculateImageBase64FromItems(app *App, items []models.ItemWithUpd, idx int, bsgItems map[string]models.BSGItem) string {
	imageHash := cache_presets.GetItemHash(items[idx], items, bsgItems)
	imageBase64, err := loadImage(app, imageHash)
	var ImageBase64 string
	if err != nil {
		ImageBase64 = ""
	} else {
		ImageBase64 = imageBase64
	}
	return ImageBase64
}

func addUIPropertiesToInventoryItems(app *App, parentId string, inventoryItems *[]models.ItemWithUpd) {
	bsgItems := app.ctx.Value(contextAllBSGItems).(map[string]models.BSGItem)
	allItems := app.ctx.Value(contextAllItems).(*models.AllItems)

	for i := range *inventoryItems {
		inventoryItem := &(*inventoryItems)[i]

		if inventoryItem.ParentID != nil && *inventoryItem.ParentID != parentId {
			continue
		}

		imageHash := cache_presets.GetItemHash(*inventoryItem, *inventoryItems, bsgItems)
		imageBase64, err := loadImage(app, imageHash)
		var ImageBase64 string
		if err != nil {
			ImageBase64 = ""
		} else {
			ImageBase64 = imageBase64
		}
		inventoryItem.ImageBase64 = ImageBase64

		sizeX, sizeY := images.GetItemSize(*inventoryItem, *inventoryItems, bsgItems)
		inventoryItem.SizeX = sizeX
		inventoryItem.SizeY = sizeY
		inventoryItem.ShortName = allItems.Items[inventoryItem.Tpl].ShortName
		inventoryItem.Amount = 1
		if inventoryItem.Upd != nil {
			inventoryItem.Amount = inventoryItem.Upd.StackObjectsCount
		}
		bsgItem, ok := bsgItems[inventoryItem.Tpl]
		if ok {
			inventoryItem.BackgroundColor = bsgItem.Props.BackgroundColor
			inventoryItem.IsStockable = bsgItem.Props.StackMaxSize != 1
		}
	}
}

func addImageToWeaponBuildAttachments(app *App, weaponBuild *models.WeaponBuild) {
	bsgItems := app.ctx.Value(contextAllBSGItems).(map[string]models.BSGItem)

	for j := range *weaponBuild.Items {
		weaponAttachment := &(*weaponBuild.Items)[j]
		bsgItem, found := bsgItems[weaponAttachment.Tpl]
		if !found {
			runtime.LogWarning(app.ctx, "Couldn't find item in BSG items: "+weaponAttachment.Tpl)
			runtime.EventsEmit(app.ctx, "toast.error", "Some attachments couldn't be loaded.")
			continue
		}
		attachmentHash := cache.GetItemHash(bsgItem, bsgItems)
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

func getCurrentActiveQuests(app *App, profile models.SPTProfile) []models.ViewQuest {
	allBsgQuests := app.ctx.Value(contextAllQuests).(*[]models.BsgQuest)
	locales := app.ctx.Value(contextLocales).(*models.Locales)

	var quests []models.ViewQuest
	for _, quest := range profile.Characters.PMC.Quests {
		if quest.Status != 2 && quest.Status != 3 {
			continue
		}

		idx := slices.IndexFunc(*allBsgQuests, func(i models.BsgQuest) bool {
			return i.Id == quest.QID
		})
		if idx == -1 {
			runtime.LogWarning(app.ctx, "Couldn't find quest with id "+quest.QID)
			continue
		}
		bsgQuest := (*allBsgQuests)[idx]

		location := bsgQuest.Location
		if location != "any" {
			location = locales.Data[fmt.Sprintf("%s Name", bsgQuest.Location)]
		}

		questItem := models.ViewQuest{
			QID:      quest.QID,
			Name:     locales.Data[bsgQuest.Name],
			Location: location,
			Trader:   locales.Data[fmt.Sprintf("%s Nickname", bsgQuest.TraderId)],
		}
		quests = append(quests, questItem)
	}
	return quests
}

func finishQuest(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		sessionId := app.ctx.Value(contextSessionId).(string)
		questId := r.FormValue("id")
		e := api.FinishQuest(app.config.GetSptUrl(), sessionId, questId)
		if e != nil {
			runtime.LogWarning(app.ctx, "Couldn't finish quest: "+e.Error())
		}
		runtime.EventsEmit(app.ctx, "toast.info", "Quest finished "+questId)
	}
}

func refreshQuest(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		err := reloadProfiles(app)
		if err != nil {
			redirectToErrorPage(app, err.Error())
			return
		}
		profile := getProfileFromSession(app)
		quests := getCurrentActiveQuests(app, profile)
		templ.Handler(components.Quests(&profile, quests)).ServeHTTP(w, r)
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
	r.Get("/stash", getStash(app))
	r.Post("/stash", addStashItem(app))
	r.Get("/kit/{id}", getKit(app))
	r.Get("/kits", getKits(app))
	r.Post("/kit/{presetId}", addKit(app))
	r.Post("/spt", sendSptMessage(app))
	// forward calls to SPT server for files (images)
	r.Get("/file", getFile(app))
	r.Get("/linked-search/{id}", getLinkedSearchModal(app))
	r.Get("/reload-profiles", goToProfileList(app))
	r.Post("/quest", finishQuest(app))
	r.Get("/quest", refreshQuest(app))
	// this is not used as it is disabled in the template
	// https://github.com/angel-git/give-ui/issues/49
	r.Post("/magazine-loadouts/{id}", addMagazineLoadout(app))

	return r
}
