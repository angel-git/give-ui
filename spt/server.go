package spt

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ServerInfo struct {
	Version    string `json:"version"`
	Path       string `json:"path"`
	ModVersion string `json:"modVersion"`
}

type SPTProfileInfo struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type SPTProfile struct {
	Info SPTProfileInfo `json:"info"`
}

type BSGItem struct {
	Id     string       `json:"_id"`
	Parent string       `json:"_parent"`
	Props  BSGItemProps `json:"_props"`
}

type BSGItemProps struct {
	StackMaxSize int `json:"StackMaxSize"`
}

type Locales struct {
	Data map[string]string `json:"data"`
}

type ViewItem struct {
	Id          string
	Name        string
	Description string
	MaxStock    int
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func doGet(url string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("responsecompressed", "0")
	return myClient.Do(req)
}

func getJson(url string, target interface{}) error {
	r, err := doGet(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func getRawBytes(url string) ([]byte, error) {
	r, err := doGet(url)
	if err != nil {
		return nil, err
	}
	// Read the entire response body into a byte slice
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func ConnectToSptServer(host string, port string) (r *ServerInfo, e error) {
	serverInfo := &ServerInfo{}
	err := getJson(fmt.Sprintf("http://%s:%s/tarkov-stash/server", host, port), serverInfo)
	if err != nil {
		return nil, err
	}
	return serverInfo, nil
}

func LoadProfiles(host string, port string) (r []SPTProfileInfo, e error) {
	profiles, err := getRawBytes(fmt.Sprintf("http://%s:%s/tarkov-stash/profiles", host, port))
	if err != nil {
		return nil, err
	}
	var sessionsMap map[string]SPTProfile
	err = parseByteResponse(profiles, &sessionsMap)
	if err != nil {
		return nil, err
	}
	fmt.Println("sessionsMap", sessionsMap)
	var sessions []SPTProfileInfo
	for _, v := range sessionsMap {
		sessions = append(sessions, v.Info)
	}
	return sessions, nil
}

func LoadItems(host string, port string) (r map[string]ViewItem, e error) {
	itemsBytes, err := getRawBytes(fmt.Sprintf("http://%s:%s/tarkov-stash/items", host, port))
	if err != nil {
		return nil, err
	}
	var itemsMap map[string]BSGItem
	err = parseByteResponse(itemsBytes, &itemsMap)
	if err != nil {
		return nil, err
	}
	// TODO hardcoded locale en
	localeBytes, err := getRawBytes(fmt.Sprintf("http://%s:%s/client/locale/en", host, port))
	if err != nil {
		return nil, err
	}
	var locales Locales
	err = parseByteResponse(localeBytes, &locales)
	if err != nil {
		return nil, err
	}

	viewItemsMap := make(map[string]ViewItem)
	for key, bsgItem := range itemsMap {
		name := locales.Data[fmt.Sprintf("%s Name", bsgItem.Id)]
		description := locales.Data[fmt.Sprintf("%s Description", bsgItem.Id)]
		viewItem := ViewItem{
			Id:          bsgItem.Id,
			Name:        name,
			Description: description,
			MaxStock:    bsgItem.Props.StackMaxSize,
		}
		viewItemsMap[key] = viewItem

	}
	return viewItemsMap, nil
}

func LoadPresets() {
	// TODO maybe we should merge both into LoadItems
}

func parseByteResponse(profiles []byte, target interface{}) error {
	err := json.Unmarshal(profiles, target)
	if err != nil {
		return err
	}
	return nil
}
