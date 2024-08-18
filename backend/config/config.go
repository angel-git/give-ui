package config

import (
	"spt-give-ui/backend/database"
	"spt-give-ui/backend/logger"
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
	defaultJsonConfig := database.JsonDatabase{
		Locale: "English",
		Theme:  defaultTheme,
		SptUrl: "http://127.0.0.1:6969",
	}
	jsonConfig := database.CreateDatabase(defaultJsonConfig)
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
	database.SaveValue(database.LocaleDbKey, locale)
}

func (c *Config) SwitchTheme() {
	newTheme := ""
	if c.theme == defaultTheme {
		newTheme = lightTheme
	} else {
		newTheme = defaultTheme
	}
	c.theme = newTheme
	database.SaveValue(database.ThemeDbKey, newTheme)
}

func (c *Config) SetSptUrl(url string) {
	c.sptUrl = url
	database.SaveValue(database.SptSeverDbKey, url)
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
