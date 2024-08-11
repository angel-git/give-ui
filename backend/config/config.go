package config

import (
	"go.etcd.io/bbolt"
	"spt-give-ui/backend/database"
	"spt-give-ui/backend/logger"
)

const (
	defaultTheme  = "mytheme"
	lightTheme    = "retro"
	localeDbKey   = "locale"
	themeDbKey    = "theme"
	sptSeverDbKey = "sptUrl"
)

type Config struct {
	errorLogger *logger.ErrorFileLogger
	db          *bbolt.DB
	locale      string
	theme       string
	sptUrl      string
}

func LoadConfig() *Config {
	errorLogger := logger.SetupLogger()
	db := database.CreateDatabase()
	locale := database.GetValue(db, localeDbKey)
	if locale == "" {
		locale = "English"
		database.SaveValue(db, localeDbKey, locale)
	}

	theme := database.GetValue(db, themeDbKey)
	if theme == "" {
		theme = defaultTheme
		database.SaveValue(db, theme, theme)
	}

	sptUrl := database.GetValue(db, sptSeverDbKey)
	if sptUrl == "" {
		sptUrl = "http://127.0.0.1:6969"
		database.SaveValue(db, sptSeverDbKey, sptUrl)
	}
	return &Config{
		errorLogger: errorLogger,
		db:          db,
		locale:      locale,
		theme:       theme,
		sptUrl:      sptUrl,
	}
}

func (c *Config) SetLocale(locale string) {
	c.locale = locale
	database.SaveValue(c.db, localeDbKey, locale)
}

func (c *Config) SwitchTheme() {
	newTheme := ""
	if c.theme == defaultTheme {
		newTheme = lightTheme
	} else {
		newTheme = defaultTheme
	}
	c.theme = newTheme
	database.SaveValue(c.db, themeDbKey, newTheme)
}

func (c *Config) SetSptUrl(url string) {
	c.sptUrl = url
	database.SaveValue(c.db, sptSeverDbKey, url)
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
	err := c.db.Close()
	if err != nil {
		return err
	}
	return c.errorLogger.Close()
}
