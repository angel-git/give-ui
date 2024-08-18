package config

import (
	"spt-give-ui/backend/logger"
	"spt-give-ui/backend/store"
)

const (
	defaultTheme = "mytheme"
	lightTheme   = "retro"
)

type Config struct {
	errorLogger *logger.ErrorFileLogger
	locale      string
	theme       string
	sptUrl      string
}

func LoadConfig() *Config {
	errorLogger := logger.SetupLogger()
	defaultJsonConfig := store.JsonDatabase{
		Locale: "English",
		Theme:  defaultTheme,
		SptUrl: "http://127.0.0.1:6969",
	}
	jsonConfig := store.CreateDatabase(defaultJsonConfig)
	return &Config{
		errorLogger: errorLogger,
		//db:          db,
		locale: jsonConfig.Locale,
		theme:  jsonConfig.Theme,
		sptUrl: jsonConfig.SptUrl,
	}
}

func (c *Config) SetLocale(locale string) {
	c.locale = locale
	store.SaveValue(store.LocaleDbKey, locale)
}

func (c *Config) SwitchTheme() {
	newTheme := ""
	if c.theme == defaultTheme {
		newTheme = lightTheme
	} else {
		newTheme = defaultTheme
	}
	c.theme = newTheme
	store.SaveValue(store.ThemeDbKey, newTheme)
}

func (c *Config) SetSptUrl(url string) {
	c.sptUrl = url
	store.SaveValue(store.SptSeverDbKey, url)
}

func (c *Config) GetLocale() string {
	return c.locale
}

func (c *Config) GetTheme() string {
	return c.theme
}

func (c *Config) GetSptUrl() string {
	return c.sptUrl
}

func (c *Config) Close() error {
	//err := c.db.Close()
	//if err != nil {
	//	return err
	//}
	return c.errorLogger.Close()
}
