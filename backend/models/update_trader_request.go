package models

type UpdateTraderRequest struct {
	Nickname string `json:"nickname"`
	Rep      string `json:"rep"`
	Spend    string `json:"spend"`
}
