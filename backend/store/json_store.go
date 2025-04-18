package store

import (
	"encoding/json"
	"log"
	"os"
)

type JsonDatabase struct {
	Locale        string   `json:"locale"`
	Theme         string   `json:"theme"`
	SptUrl        string   `json:"sptUrl"`
	CacheFolder   string   `json:"cacheFolder"`
	FavoriteItems []string `json:"favoriteItems"`
	IgnoreCache   bool     `json:"ignoreCache"`
}

const LocaleDbKey = "locale"
const ThemeDbKey = "theme"
const SptSeverDbKey = "sptUrl"
const FavoriteItemsDbKey = "favoriteItems"
const CacheFolderDbKey = "cacheFolder"
const IgnoreCacheDbKey = "ignoreCache"

const dbName = "give-ui.config.json"

func CreateDatabase(defaultConfig JsonDatabase) JsonDatabase {
	file, _ := os.Open(dbName)
	if file == nil {
		content, _ := json.Marshal(defaultConfig)
		err := os.WriteFile(dbName, content, 0600)
		if err != nil {
			log.Fatalf("Error creating json config file: %s", err)
		}
		return defaultConfig
	}
	defer file.Close()

	content, err := os.ReadFile(dbName)
	if err != nil {
		log.Fatalf("Error reading json config file: %s", err)
	}
	jsonConfig := JsonDatabase{}
	err = json.Unmarshal(content, &jsonConfig)
	if err != nil {
		log.Fatalf("Error reading json config content: %s", err)
	}

	return jsonConfig

}

func SaveValue(key string, value any) {
	content, err := os.ReadFile(dbName)
	if err != nil {
		log.Fatalf("Error writing key [%s] with value [%s]: %s", key, value, err)
	}
	jsonConfig := map[string]any{}
	err = json.Unmarshal(content, &jsonConfig)
	if err != nil {
		log.Fatalf("Error writing key [%s] with value [%s]: %s", key, value, err)
	}
	jsonConfig[key] = value
	newContent, err := json.Marshal(jsonConfig)
	err = os.WriteFile(dbName, newContent, 0600)
	if err != nil {
		log.Fatalf("Error writing key [%s] with value [%s]: %s", key, value, err)
	}

}
