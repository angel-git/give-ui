package models

type ServerInfo struct {
	Version    string `json:"version"`
	Path       string `json:"path"`
	ModVersion string `json:"modVersion"`
}
