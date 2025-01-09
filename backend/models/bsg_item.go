package models

type GlobalPreset struct {
	Id string `json:"_id"`
	//_type: string
	//_name: string
	//_parent: string
	Items []ItemWithUpd `json:"_items"`
	/** Default presets have this property */
	Encyclopedia string `json:"_encyclopedia"`
}

type BSGItem struct {
	Id     string       `json:"_id"`
	Parent string       `json:"_parent"`
	Name   string       `json:"_name"`
	Type   string       `json:"_type"`
	Props  BSGItemProps `json:"_props"`
}

type BSGItemProps struct {
	StackMaxSize            int    `json:"StackMaxSize"`
	IsUnbuyable             bool   `json:"IsUnbuyable"`
	HasHinge                bool   `json:"HasHinge"`
	Foldable                bool   `json:"Foldable"`
	VisibleAmmoRangesString string `json:"VisibleAmmoRangesString"`
	HideEntrails            bool   `json:"HideEntrails"`
	Cartridges              *[]any `json:"Cartridges"`
}
