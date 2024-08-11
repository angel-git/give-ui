package config

import (
	"go.etcd.io/bbolt"
	"spt-give-ui/backend/database"
)

type Config struct {
	db     *bbolt.DB
	locale string
	theme  string
}

func LoadConfig() *Config {
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
		db:     db,
		locale: locale,
		theme:  theme,
	}
}

func (c *Config) SetLocale(locale string) {
	// write in database AND refresh config
	c.locale = locale
	database.SaveValue(c.db, "locale", locale)
}

func (c *Config) Close() error {
	return c.db.Close()
}

func (c *Config) GetLocale() string {
	return c.locale
}
