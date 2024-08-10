package models

type AddItemRequest struct {
	ItemId string `json:"itemId"`
	Amount int    `json:"amount"`
}

type AddUserWeaponPresetRequest struct {
	ItemId string `json:"itemId"`
}

type AddGearPresetRequest struct {
	PresetId string `json:"presetId"`
	ItemId   string `json:"itemId"`
}
