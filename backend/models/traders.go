package models

type AllTradersResponse struct {
	Traders []TraderResponse `json:"data"`
}

type TraderResponse struct {
	Id              string `json:"_id"`
	Nickname        string `json:"nickname"`
	Avatar          string `json:"avatar"`
	AvailableInRaid bool   `json:"availableInRaid"`
}

type Trader struct {
	Id             string `json:"_id"`
	Nickname       string `json:"nickname"`
	NicknameLocale string `json:"nicknameLocale"`
	Image          string `json:"avatar"`
	SalesSum       string `json:"salesSum"`
	Reputation     string `json:"standing"`
	LoyaltyLevel   int    `json:"loyaltyLevel"`
}
