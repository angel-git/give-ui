package main

import (
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
	}
	a.localeMenu = a.menu.AddSubmenu("Locale")
	a.localeMenu.AddRadio("English", true, nil, a.setLocale)
	a.localeMenu.AddRadio("Czech", false, nil, a.setLocale)
	a.localeMenu.AddRadio("French", false, nil, a.setLocale)
	a.localeMenu.AddRadio("German", false, nil, a.setLocale)
	a.localeMenu.AddRadio("Hungarian", false, nil, a.setLocale)
	a.localeMenu.AddRadio("Italian", false, nil, a.setLocale)
	a.localeMenu.AddRadio("Japanese", false, nil, a.setLocale)
	a.localeMenu.AddRadio("Korean", false, nil, a.setLocale)
	a.localeMenu.AddRadio("Polish", false, nil, a.setLocale)
	a.localeMenu.AddRadio("Portuguese", false, nil, a.setLocale)
	a.localeMenu.AddRadio("Slovak", false, nil, a.setLocale)
	a.localeMenu.AddRadio("Spanish", false, nil, a.setLocale)
	a.localeMenu.AddRadio("Spanish - Mexico", false, nil, a.setLocale)
	a.localeMenu.AddRadio("Turkish", false, nil, a.setLocale)
	a.localeMenu.AddRadio("Romanian", false, nil, a.setLocale)
	a.localeMenu.AddRadio("Русский", false, nil, a.setLocale)
}
