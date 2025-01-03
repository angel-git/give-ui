package main

import (
	"context"
	"embed"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	runtimeWails "github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"net/http"
	"runtime"
)

//go:embed all:frontend/dist components
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

//go:embed wails.json
var wailsJson string

func main() {
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
		Menu:             app.menu,
		Logger:           nil,
		LogLevel:         logger.DEBUG,
		OnStartup:        app.startup,
		OnDomReady:       app.domReady,
		OnBeforeClose:    app.beforeClose,
		OnShutdown:       app.shutdown,
		WindowStartState: options.Normal,
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
		log.Fatal(err)
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
