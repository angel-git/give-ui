package models

type ItemsResponse struct {
	Items         map[string]BSGItem      `json:"items"`
	GlobalPresets map[string]GlobalPreset `json:"globalPresets"`
}
