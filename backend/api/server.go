package api

import (
	"fmt"
	"slices"
	"sort"
	"spt-give-ui/backend/http"
	"spt-give-ui/backend/models"
	"spt-give-ui/backend/util"
	"strings"
)

func ConnectToSptServer(host string, port string) (r *models.ServerInfo, e error) {
	serverInfo := &models.ServerInfo{}
	err := util.GetJson(fmt.Sprintf("http://%s:%s/give-ui/server", host, port), serverInfo)
	if err != nil {
		return nil, err
	}
	return serverInfo, nil
}

func LoadProfiles(host string, port string) (r []models.SPTProfileInfo, e error) {
	profiles, err := util.GetRawBytes(fmt.Sprintf("http://%s:%s/give-ui/profiles", host, port))
	if err != nil {
		return nil, err
	}
	var sessionsMap map[string]models.SPTProfile
	err = util.ParseByteResponse(profiles, &sessionsMap)
	if err != nil {
		return nil, err
	}
	var sessions []models.SPTProfileInfo
	for _, v := range sessionsMap {
		sessions = append(sessions, v.Info)
	}
	return sessions, nil
}

func LoadItems(host string, port string) (r *models.AllItems, e error) {
	itemsBytes, err := util.GetRawBytes(fmt.Sprintf("http://%s:%s/give-ui/items", host, port))
	if err != nil {
		return nil, err
	}
	var itemsMap map[string]models.BSGItem
	err = util.ParseByteResponse(itemsBytes, &itemsMap)
	if err != nil {
		return nil, err
	}
	// TODO hardcoded locale en
	localeBytes, err := util.GetRawBytes(fmt.Sprintf("http://%s:%s/client/locale/en", host, port))
	if err != nil {
		return nil, err
	}
	var locales models.Locales
	err = util.ParseByteResponse(localeBytes, &locales)
	if err != nil {
		return nil, err
	}

	allItems := models.AllItems{
		Categories: []string{},
		Items:      []models.ViewItem{},
		Presets:    []models.ViewItem{},
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

		viewItem := models.ViewItem{
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

func AddItem(host string, port string, sessionId string, itemId string, amount int) (e error) {
	request := models.AddItemRequest{
		ItemId: itemId,
		Amount: amount,
	}
	_, err := http.DoPost(fmt.Sprintf("http://%s:%s/give-ui/give", host, port), sessionId, request)
	return err
}
