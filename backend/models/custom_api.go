package models

type ItemsResponse struct {
	Items         map[string]BSGItem      `json:"items"`
	GlobalPresets map[string]GlobalPreset `json:"globalPresets"`
}

type ItemsRawResponse struct {
	Items map[string]interface{} `json:"items"`
}
