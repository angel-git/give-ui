package config

import (
	"go.etcd.io/bbolt"
	"spt-give-ui/backend/database"
	"spt-give-ui/backend/logger"
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
	locale := database.GetValue(db, "locale")
	if locale == "" {
		locale = "English"
		database.SaveValue(db, "locale", locale)
	}

	theme := database.GetValue(db, "theme")
	if theme == "" {
		theme = "mytheme"
		database.SaveValue(db, "theme", theme)
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
	database.SaveValue(c.db, "locale", locale)
}

func (c *Config) Close() error {
	err := c.db.Close()
	if err != nil {
		return err
	}
	return c.errorLogger.Close()
}

func (c *Config) GetLocale() string {
	return c.locale
}
