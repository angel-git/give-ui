package models

type ViewItem struct {
	Id          string
	Name        string
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
}

type AllItems struct {
	Categories    []string
	Items         map[string]ViewItem
	GlobalPresets []ViewPreset
}

type ViewWeaponBuild struct {
	Id          string
	Name        string
	ImageBase64 string
	Items       []WeaponBuildItem
}
