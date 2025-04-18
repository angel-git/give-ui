package models

type ViewItem struct {
	Id          string
	Name        string
	ShortName   string
	Type        string
	Description string
	Category    string
	MaxStock    int
	Favorite    bool
	ImageBase64 string
}

type ViewPreset struct {
	Id           string
	Encyclopedia string
	Items        []ItemWithUpd
}

type AllItems struct {
	Categories    []string
	Items         map[string]ViewItem
	GlobalPresets []ViewPreset
}
