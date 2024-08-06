package models

type AddItemRequest struct {
	ItemId string `json:"itemId"`
	Amount int    `json:"amount"`
}

type AddUserWeaponPresetRequest struct {
	ItemId string `json:"itemId"`
}
