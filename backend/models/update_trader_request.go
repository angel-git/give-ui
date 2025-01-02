package models

type UpdateTraderRepRequest struct {
	Nickname string `json:"nickname"`
	Rep      string `json:"rep"`
}

type UpdateTraderSpendRequest struct {
	Nickname string `json:"nickname"`
	Spend    string `json:"spend"`
}
