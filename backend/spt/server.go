package spt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"sort"
	"strings"
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
	Type   string       `json:"_type"`
	Props  BSGItemProps `json:"_props"`
}

type BSGItemProps struct {
	StackMaxSize int  `json:"StackMaxSize"`
	IsUnbuyable  bool `json:"IsUnbuyable"`
}

type Locales struct {
	Data map[string]string `json:"data"`
}

type ViewItem struct {
	Id          string
	Name        string
	Type        string
	Description string
	Category    string
	MaxStock    int
}

type AllItems struct {
	Categories []string
	Items      []ViewItem
	Presets    []ViewItem
}

type AddItemRequest struct {
	ItemId string `json:"itemId"`
	Amount int    `json:"amount"`
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func doPost(url string, sessionId string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("responsecompressed", "0")
	req.Header.Set("requestcompressed", "0")
	req.Header.Set("Cookie", fmt.Sprintf("PHPSESSID=%s", sessionId))
	return myClient.Do(req)
}

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
	err := getJson(fmt.Sprintf("http://%s:%s/give-ui/server", host, port), serverInfo)
	if err != nil {
		return nil, err
	}
	return serverInfo, nil
}

func LoadProfiles(host string, port string) (r []SPTProfileInfo, e error) {
	profiles, err := getRawBytes(fmt.Sprintf("http://%s:%s/give-ui/profiles", host, port))
	if err != nil {
		return nil, err
	}
	var sessionsMap map[string]SPTProfile
	err = parseByteResponse(profiles, &sessionsMap)
	if err != nil {
		return nil, err
	}
	var sessions []SPTProfileInfo
	for _, v := range sessionsMap {
		sessions = append(sessions, v.Info)
	}
	return sessions, nil
}

func LoadItems(host string, port string) (r *AllItems, e error) {
	itemsBytes, err := getRawBytes(fmt.Sprintf("http://%s:%s/give-ui/items", host, port))
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

	allItems := AllItems{
		Categories: []string{},
		Items:      []ViewItem{},
		Presets:    []ViewItem{},
	}
	for _, bsgItem := range itemsMap {
		if bsgItem.Type == "Node" || bsgItem.Props.IsUnbuyable {
			continue
		}
		var category string
		var parent = locales.Data[fmt.Sprintf("%s Name", bsgItem.Parent)]
		var parentParent = locales.Data[fmt.Sprintf("%s Name", itemsMap[bsgItem.Parent].Parent)]
		if parent != "" {
			category = parent
		} else if parentParent != "" {
			category = parentParent
		} else {
			continue
		}
		// filter out useless categories
		if strings.Contains(category, "Stash") ||
			strings.Contains(category, "Searchable item") ||
			strings.Contains(category, "Compound item") ||
			strings.Contains(category, "Loot container") ||
			strings.Contains(category, "Inventory") {
			continue
		}
		name := locales.Data[fmt.Sprintf("%s Name", bsgItem.Id)]
		description := locales.Data[fmt.Sprintf("%s Description", bsgItem.Id)]
		// filter out useless items
		if strings.Contains(name, "DO_NOT_USE") || name == "" {
			continue
		}

		viewItem := ViewItem{
			Id:          bsgItem.Id,
			Name:        name,
			Type:        bsgItem.Type,
			Description: description,
			Category:    category,
			MaxStock:    bsgItem.Props.StackMaxSize,
		}
		allItems.Items = append(allItems.Items, viewItem)
		if !slices.Contains(allItems.Categories, category) {
			allItems.Categories = append(allItems.Categories, category)
		}
	}
	sort.Strings(allItems.Categories)
	sort.SliceStable(allItems.Items, func(i, j int) bool {
		return allItems.Items[i].Name < allItems.Items[j].Name
	})

	return &allItems, nil
}

func LoadPresets() {
	// TODO maybe we should merge both into LoadItems
}

func AddItem(host string, port string, sessionId string, itemId string, amount int) {
	request := AddItemRequest{
		ItemId: itemId,
		Amount: amount,
	}
	doPost(fmt.Sprintf("http://%s:%s/give-ui/give", host, port), sessionId, request)
}

func parseByteResponse(profiles []byte, target interface{}) error {
	err := json.Unmarshal(profiles, target)
	if err != nil {
		return err
	}
	return nil
}
