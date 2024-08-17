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

func ConnectToSptServer(url string) (r *models.ServerInfo, e error) {
	serverInfo := &models.ServerInfo{}
	err := util.GetJson(fmt.Sprintf("%s/give-ui/server", url), serverInfo)
	if err != nil {
		return nil, err
	}
	return serverInfo, nil
}

func LoadProfiles(url string) (r []models.SPTProfile, e error) {
	profiles, err := util.GetRawBytes(fmt.Sprintf("%s/give-ui/profiles", url))
	if err != nil {
		return nil, err
	}
	var sessionsMap map[string]models.SPTProfile
	err = util.ParseByteResponse(profiles, &sessionsMap)
	if err != nil {
		return nil, err
	}
	var sessions []models.SPTProfile
	for _, v := range sessionsMap {
		sessions = append(sessions, v)
	}
	sort.SliceStable(sessions, func(i, j int) bool {
		return sessions[i].Info.Username < sessions[j].Info.Username
	})
	return sessions, nil
}

func LoadItems(url string, locale string) (r *models.AllItems, e error) {
	items, err := getItemsFromServer(url)
	if err != nil {
		return nil, err
	}
	locales, err := getLocaleFromServer(url, locale)
	if err != nil {
		return nil, err
	}

	allItems := parseItems(items, *locales)

	return &allItems, nil
}

func AddItem(url string, sessionId string, itemId string, amount int) (e error) {
	request := models.AddItemRequest{
		ItemId: itemId,
		Amount: amount,
	}
	_, err := http.DoPost(fmt.Sprintf("%s/give-ui/give", url), sessionId, request)
	return err
}

func AddUserWeapon(url string, sessionId string, presetId string) (e error) {
	request := models.AddUserWeaponPresetRequest{
		ItemId: presetId,
	}
	_, err := http.DoPost(fmt.Sprintf("%s/give-ui/give-user-preset", url), sessionId, request)
	return err
}

func getLocaleFromServer(url string, locale string) (*models.Locales, error) {
	localeBytes, err := util.GetRawBytes(fmt.Sprintf("%s/client/locale/%s", url, locale))
	if err != nil {
		return nil, err
	}
	var locales *models.Locales
	err = util.ParseByteResponse(localeBytes, &locales)
	if err != nil {
		return nil, err
	}
	return locales, nil
}

func getItemsFromServer(url string) (*models.ItemsResponse, error) {
	itemsBytes, err := util.GetRawBytes(fmt.Sprintf("%s/give-ui/items", url))
	if err != nil {
		return nil, err
	}
	var itemsMap *models.ItemsResponse
	err = util.ParseByteResponse(itemsBytes, &itemsMap)
	if err != nil {
		return nil, err
	}
	return itemsMap, nil
}

func parseItems(items *models.ItemsResponse, locales models.Locales) models.AllItems {
	const NameFormat = "%s Name"
	const DescriptionFormat = "%s Description"
	allItems := models.AllItems{
		Categories:    []string{},
		Items:         map[string]models.ViewItem{},
		GlobalPresets: []models.ViewPreset{},
	}

	for _, globalPreset := range items.GlobalPresets {
		viewPreset := models.ViewPreset{
			Id:           globalPreset.Id,
			Encyclopedia: globalPreset.Encyclopedia,
		}
		allItems.GlobalPresets = append(allItems.GlobalPresets, viewPreset)
	}

	itemsMap := items.Items
	for _, bsgItem := range itemsMap {
		if bsgItem.Type == "Node" || bsgItem.Props.IsUnbuyable {
			continue
		}
		// filter test broken items
		if slices.Contains(getHiddenItems(), bsgItem.Id) {
			continue
		}

		var category string
		var parent = locales.Data[fmt.Sprintf(NameFormat, bsgItem.Parent)]
		var parentParent = locales.Data[fmt.Sprintf(NameFormat, itemsMap[bsgItem.Parent].Parent)]
		if parent != "" {
			category = parent
		} else if parentParent != "" {
			category = parentParent
		} else {
			continue
		}
		// filter out useless categories
		if slices.Contains(getHiddenCategories(), bsgItem.Parent) {
			continue
		}
		name := locales.Data[fmt.Sprintf(NameFormat, bsgItem.Id)]
		description := locales.Data[fmt.Sprintf(DescriptionFormat, bsgItem.Id)]
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
		allItems.Items[viewItem.Id] = viewItem
		if !slices.Contains(allItems.Categories, category) {
			allItems.Categories = append(allItems.Categories, category)
		}
	}
	sort.Strings(allItems.Categories)
	return allItems
}

func getHiddenCategories() []string {
	return []string{
		"55d720f24bdc2d88028b456d",
		"62f109593b54472778797866",
		"63da6da4784a55176c018dba",
		"566abbb64bdc2d144c8b457d",
		"566965d44bdc2d814c8b4571",
	}
}

func getHiddenItems() []string {
	return []string{
		"5ae083b25acfc4001a5fc702",
	}
}
