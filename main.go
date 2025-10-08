package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"runtime"
	giveLogger "spt-give-ui/backend/logger"

	"github.com/tidwall/gjson"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	runtimeWails "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist components
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

//go:embed wails.json
var wailsJson string

func main() {
	log := giveLogger.SetupLogger()
	version := gjson.Get(wailsJson, "version").Str
	name := gjson.Get(wailsJson, "name").Str
	// Create an instance of the app structure and custom Middleware
	app := NewApp(name, version)
	app.makeMenu()
	r := NewChiRouter(app)

	// Create application with options
	err := wails.Run(&options.App{
		Title:             "give-ui",
		Width:             1000,
		Height:            700,
		MinWidth:          1000,
		MinHeight:         700,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		AssetServer: &assetserver.Options{
			Assets: assets,
			Middleware: func(next http.Handler) http.Handler {
				r.NotFound(next.ServeHTTP)
				return r
			},
		},
		Menu:               app.menu,
		Logger:             log,
		LogLevel:           logger.INFO,
		LogLevelProduction: logger.INFO,
		OnStartup:          app.startup,
		OnDomReady:         app.domReady,
		OnBeforeClose:      app.beforeClose,
		OnShutdown:         app.shutdown,
		WindowStartState:   options.Normal,
		Bind: []interface{}{
			app,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
			// DisableFramelessWindowDecorations: false,
			WebviewUserDataPath: "",
			ZoomFactor:          1.0,
		},
		// Mac platform specific options
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: false,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "give-ui",
				Message: fmt.Sprintf("Version: %s", version),
				Icon:    icon,
			},
		},
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}

func (a *App) makeMenu() {
	a.menu = menu.NewMenu()
	if runtime.GOOS == "darwin" {
		a.menu.Append(menu.AppMenu())
		a.menu.Append(menu.EditMenu())
	}
	localeFromConfig := a.config.GetLocale()
	a.localeMenu = a.menu.AddSubmenu("Locale")
	a.localeMenu.Append(addRadio("English", localeFromConfig, a.setLocale))
	a.localeMenu.Append(addRadio("Chinese", localeFromConfig, a.setLocale))
	a.localeMenu.Append(addRadio("Czech", localeFromConfig, a.setLocale))
	a.localeMenu.Append(addRadio("French", localeFromConfig, a.setLocale))
	a.localeMenu.Append(addRadio("German", localeFromConfig, a.setLocale))
	a.localeMenu.Append(addRadio("Hungarian", localeFromConfig, a.setLocale))
	a.localeMenu.Append(addRadio("Italian", localeFromConfig, a.setLocale))
	a.localeMenu.Append(addRadio("Japanese", localeFromConfig, a.setLocale))
	a.localeMenu.Append(addRadio("Korean", localeFromConfig, a.setLocale))
	a.localeMenu.Append(addRadio("Polish", localeFromConfig, a.setLocale))
	a.localeMenu.Append(addRadio("Portuguese", localeFromConfig, a.setLocale))
	a.localeMenu.Append(addRadio("Slovak", localeFromConfig, a.setLocale))
	a.localeMenu.Append(addRadio("Spanish", localeFromConfig, a.setLocale))
	a.localeMenu.Append(addRadio("Spanish - Mexico", localeFromConfig, a.setLocale))
	a.localeMenu.Append(addRadio("Turkish", localeFromConfig, a.setLocale))
	a.localeMenu.Append(addRadio("Romanian", localeFromConfig, a.setLocale))
	a.localeMenu.Append(addRadio("Русский", localeFromConfig, a.setLocale))

	a.settingsMenu = a.menu.AddSubmenu("Settings")
	if a.config.GetCacheFolder() == "" {
		a.settingsMenu.Append(menu.Text("Select cache folder", nil, a.selectCacheFolder))
	} else {
		a.settingsMenu.Append(menu.Text("Use default cache folder", nil, a.clearCacheFolder))
	}
	a.settingsMenu.Append(menu.Checkbox("Load images from cache", a.config.GetUseCache(), nil, a.toggleUseCache))
	a.settingsMenu.Append(menu.Checkbox("Log profiles response", a.config.GetLogResponses(), nil, a.toggleLogResponses))
	a.settingsMenu.Append(menu.Text("Switch server", nil, a.switchServer))
}

func addRadio(label string, selected string, click menu.Callback) *menu.MenuItem {
	item := menu.Radio(label, label == selected, nil, click)
	return item
}

func (a *App) setLocale(data *menu.CallbackData) {
	if a.config.GetLocale() == data.MenuItem.Label {
		return
	}
	a.config.SetLocale(data.MenuItem.Label)
	a.ctx = context.WithValue(a.ctx, contextLocales, nil)
	for _, localeMenu := range a.localeMenu.Items {
		localeMenu.Checked = false
	}
	data.MenuItem.Checked = true

	// refresh menu with the selected locale
	runtimeWails.MenuSetApplicationMenu(a.ctx, a.menu)
	runtimeWails.MenuUpdateApplicationMenu(a.ctx)

	// refresh to main screen
	runtimeWails.WindowReloadApp(a.ctx)
}

func (a *App) selectCacheFolder(data *menu.CallbackData) {
	folder, err := runtimeWails.OpenDirectoryDialog(a.ctx, runtimeWails.OpenDialogOptions{
		Title: "Example c:\\games\\spt\\user\\sptappdata\\live",
	})
	if err != nil {
		runtimeWails.MessageDialog(a.ctx, runtimeWails.MessageDialogOptions{
			Type:    runtimeWails.ErrorDialog,
			Title:   "Error",
			Message: err.Error(),
		})
		return
	}
	a.config.SetCacheFolder(folder)
	data.MenuItem.Label = "Use default cache folder"
	data.MenuItem.OnClick(a.clearCacheFolder)

	// refresh menu with the selected locale

	runtimeWails.MenuSetApplicationMenu(a.ctx, a.menu)
	runtimeWails.MenuUpdateApplicationMenu(a.ctx)

	// refresh to main screen
	runtimeWails.WindowReloadApp(a.ctx)
}

func (a *App) clearCacheFolder(data *menu.CallbackData) {
	a.config.SetCacheFolder("")

	data.MenuItem.Label = "Select cache folder"
	data.MenuItem.OnClick(a.selectCacheFolder)

	// refresh menu with the selected locale
	runtimeWails.MenuSetApplicationMenu(a.ctx, a.menu)
	runtimeWails.MenuUpdateApplicationMenu(a.ctx)

	// refresh to main screen
	runtimeWails.WindowReloadApp(a.ctx)
}

func (a *App) toggleUseCache(data *menu.CallbackData) {
	a.config.SetUseCache(!a.config.GetUseCache())
	data.MenuItem.Checked = a.config.GetUseCache()

	// refresh menu with the selected locale
	runtimeWails.MenuSetApplicationMenu(a.ctx, a.menu)
	runtimeWails.MenuUpdateApplicationMenu(a.ctx)

	// refresh to main screen
	runtimeWails.WindowReloadApp(a.ctx)
}

func (a *App) toggleLogResponses(data *menu.CallbackData) {
	a.config.SetLogResponses(!a.config.GetLogResponses())
	data.MenuItem.Checked = a.config.GetLogResponses()

	// refresh menu with the selected locale
	runtimeWails.MenuSetApplicationMenu(a.ctx, a.menu)
	runtimeWails.MenuUpdateApplicationMenu(a.ctx)
	// refresh to main screen
	runtimeWails.WindowReloadApp(a.ctx)
}

func (a *App) switchServer(_data *menu.CallbackData) {
	// refresh to main screen
	runtimeWails.WindowReloadApp(a.ctx)
}
