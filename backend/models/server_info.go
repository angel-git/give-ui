package models

type ServerInfo struct {
	Version    string          `json:"version"`
	Path       string          `json:"path"`
	ModVersion string          `json:"modVersion"`
	MaxLevel   int             `json:"maxLevel"`
	Gifts      map[string]Gift `json:"gifts"`
}

type Gift struct {
	Items []ItemWithUpd `json:"items"`
}
