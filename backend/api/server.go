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
	err := util.GetJson(fmt.Sprintf("%s/give-ui/server", url), "", serverInfo)
	if err != nil {
		return nil, err
	}
	return serverInfo, nil
}

func LoadProfiles(url string) (r []models.SPTProfile, e error) {
	profiles, err := util.GetRawBytes(fmt.Sprintf("%s/give-ui/profiles", url), "")
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

func LoadItems(url string) (r *models.ItemsResponse, e error) {
	items, err := getItemsFromServer(url)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func ParseItems(allItems *models.ItemsResponse, locales *models.Locales) (r *models.AllItems, e error) {
	items := parseItems(allItems, *locales)
	return &items, nil
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

func LoadSkills(profile models.SPTProfile, locales *models.Locales) (r []models.Skill, e error) {
	var skills []models.Skill
	// try to find skill in lowercase, Troubleshooting example
	localesLowCase := convertLocalesToLowercase(locales)
	for _, skill := range profile.Characters.PMC.Skills.Common {
		name, foundName := localesLowCase[strings.ToLower(skill.Id)]
		if !foundName {
			continue
		}
		skills = append(skills, models.Skill{
			Id:       skill.Id,
			Name:     name,
			Progress: fmt.Sprintf("%d", int(skill.Progress/100)),
		})
	}
	sort.SliceStable(skills, func(i, j int) bool {
		return skills[i].Name < skills[j].Name
	})

	return skills, nil
}

func convertLocalesToLowercase(locales *models.Locales) map[string]string {
	localeLowercase := make(map[string]string)
	for k, v := range locales.Data {
		localeLowercase[strings.ToLower(k)] = v
	}
	return localeLowercase
}

func LoadTraders(url string, profile models.SPTProfile, sessionId string, locales *models.Locales) (r []models.Trader, e error) {
	tradersResponse := &models.AllTradersResponse{}
	err := util.GetJson(fmt.Sprintf("%s/client/trading/api/traderSettings", url), sessionId, tradersResponse)
	if err != nil {
		return nil, err
	}
	traders := parseTraders(url, tradersResponse, profile, locales)
	return traders, nil
}

func UpdateTraderSpend(url string, sessionId string, nickname string, spend string) (e error) {
	request := models.UpdateTraderSpendRequest{
		Nickname: nickname,
		Spend:    spend,
	}
	_, err := http.DoPost(fmt.Sprintf("%s/give-ui/update-trader-spend", url), sessionId, request)
	return err
}
func UpdateTraderRep(url string, sessionId string, nickname string, rep string) (e error) {
	request := models.UpdateTraderRepRequest{
		Nickname: nickname,
		Rep:      rep,
	}
	_, err := http.DoPost(fmt.Sprintf("%s/give-ui/update-trader-rep", url), sessionId, request)
	return err
}

func UpdateLevel(url string, sessionId string, level int) (e error) {
	request := models.UpdateLevelRequest{
		Level: level,
	}
	_, err := http.DoPost(fmt.Sprintf("%s/give-ui/update-level", url), sessionId, request)
	return err
}

func UpdateSkill(url string, sessionId string, skill string, progress int) (e error) {
	request := models.UpdateSkillRequest{
		Progress: progress,
		Skill:    skill,
	}
	_, err := http.DoPost(fmt.Sprintf("%s/give-ui/update-skill", url), sessionId, request)
	return err
}

func LoadImage(url string, sessionId string, imageHash string) (r string, e error) {
	response := &models.CacheImageResponse{}
	err := util.GetJson(fmt.Sprintf("%s/give-ui/cache/%s", url, imageHash), sessionId, response)
	if err != nil {
		return "", err
	}
	if response.Error != nil {
		return "", fmt.Errorf(*response.Error)
	}
	return *response.ImageBase64, nil
}

func parseTraders(url string, tradersResponse *models.AllTradersResponse, profile models.SPTProfile, locales *models.Locales) []models.Trader {

	var traders []models.Trader
	for _, trader := range tradersResponse.Traders {
		traderProfile, foundTrader := profile.Characters.PMC.TradersInfo[trader.Id]
		if !foundTrader || trader.AvailableInRaid {
			continue
		}
		var nicknameLocale = locales.Data[fmt.Sprintf("%s Nickname", trader.Id)]
		var maxRep string
		var loyaltyLevel = traderProfile.LoyaltyLevel
		if trader.Id == "579dc571d53a0658a154fbec" {
			// fence
			maxRep = "7"
			if loyaltyLevel == 2 {
				loyaltyLevel = 4
			}
		} else {
			maxRep = "2"
		}
		traders = append(traders, models.Trader{
			Id:             trader.Id,
			Nickname:       trader.Nickname,
			NicknameLocale: nicknameLocale,
			Reputation:     fmt.Sprintf("%.2f", traderProfile.Standing),
			SalesSum:       fmt.Sprintf("%.0f", traderProfile.SalesSum),
			Image:          fmt.Sprintf("%s%s", url, trader.Avatar),
			MaxRep:         maxRep,
			LoyaltyLevel:   loyaltyLevel,
		})
	}
	sort.SliceStable(traders, func(i, j int) bool {
		return traders[i].Id < traders[j].Id
	})

	return traders
}

func GetLocaleFromServer(url string, locale string) (*models.Locales, error) {
	localeBytes, err := util.GetRawBytes(fmt.Sprintf("%s/client/locale/%s", url, locale), "")
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
	itemsBytes, err := util.GetRawBytes(fmt.Sprintf("%s/give-ui/items", url), "")
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
			Items:        globalPreset.Items,
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
			ImageBase64: "",
			Category:    category,
			MaxStock:    bsgItem.Props.StackMaxSize,
			Favorite:    false,
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
