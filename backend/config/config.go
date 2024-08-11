package config

import (
	"go.etcd.io/bbolt"
	"spt-give-ui/backend/database"
	"spt-give-ui/backend/logger"
)

const (
	defaultTheme = "mytheme"
	lightTheme   = "retro"
	localeDbKey  = "locale"
	themeDbKey   = "theme"
)

type Config struct {
	errorLogger *logger.ErrorFileLogger
	db          *bbolt.DB
	locale      string
	theme       string
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

	return &Config{
		errorLogger: errorLogger,
		db:          db,
		locale:      locale,
		theme:       theme,
	}
}

func (c *Config) SetLocale(locale string) {
	// write in database AND refresh config
	c.locale = locale
	database.SaveValue(c.db, localeDbKey, locale)
}

func (c *Config) SwitchTheme() {
	// write in database AND refresh config
	newTheme := ""
	if c.theme == defaultTheme {
		newTheme = lightTheme
	} else {
		newTheme = defaultTheme
	}
	c.theme = newTheme
	database.SaveValue(c.db, themeDbKey, newTheme)
}

func (c *Config) GetLocale() string {
	return c.locale
}

func (c *Config) GetTheme() string {
	return c.theme
}

func (c *Config) Close() error {
	err := c.db.Close()
	if err != nil {
		return err
	}
	return c.errorLogger.Close()
}
