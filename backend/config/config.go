package config

import (
	"slices"
	"spt-give-ui/backend/logger"
	"spt-give-ui/backend/store"
)

const (
	defaultTheme = "mytheme"
	lightTheme   = "retro"
)

type Config struct {
	errorLogger   *logger.FileLogger
	locale        string
	theme         string
	sptUrl        string
	cacheFolder   string
	favoriteItems []string
	useCache      bool
}

func LoadConfig() *Config {
	errorLogger := logger.SetupLogger()
	defaultJsonConfig := store.JsonDatabase{
		Locale:        "English",
		Theme:         defaultTheme,
		SptUrl:        "https://127.0.0.1:6969",
		CacheFolder:   "",
		FavoriteItems: []string{},
		IgnoreCache:   false,
	}
	jsonConfig := store.CreateDatabase(defaultJsonConfig)
	return &Config{
		errorLogger: errorLogger,
		//db:          db,
		locale:        jsonConfig.Locale,
		theme:         jsonConfig.Theme,
		sptUrl:        jsonConfig.SptUrl,
		favoriteItems: jsonConfig.FavoriteItems,
		cacheFolder:   jsonConfig.CacheFolder,
		useCache:      !jsonConfig.IgnoreCache,
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

func (c *Config) ToggleFavoriteItem(itemId string) {
	idx := slices.Index(c.favoriteItems, itemId)
	if slices.Contains(c.favoriteItems, itemId) {
		c.favoriteItems = remove(c.favoriteItems, idx)
	} else {
		c.favoriteItems = append(c.favoriteItems, itemId)
	}
	store.SaveValue(store.FavoriteItemsDbKey, c.favoriteItems)
}

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
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

func (c *Config) GetFavoriteItems() []string {
	return c.favoriteItems
}
func (c *Config) GetCacheFolder() string {
	return c.cacheFolder
}

func (c *Config) SetCacheFolder(folder string) {
	c.cacheFolder = folder
	store.SaveValue(store.CacheFolderDbKey, folder)
}

func (c *Config) GetUseCache() bool {
	return c.useCache
}

func (c *Config) SetUseCache(cache bool) {
	c.useCache = cache
	store.SaveValue(store.IgnoreCacheDbKey, !cache)
}

func (c *Config) Close() error {
	//err := c.db.Close()
	//if err != nil {
	//	return err
	//}
	return c.errorLogger.Close()
}
