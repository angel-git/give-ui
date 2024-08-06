package models

type GlobalPreset struct {
	Id string `json:"_id"`
	//_type: string
	//_name: string
	//_parent: string
	//_items: Item[]
	/** Default presets have this property */
	Encyclopedia string `json:"_encyclopedia"`
}
type BSGItem struct {
	Id     string       `json:"_id"`
	Parent string       `json:"_parent"`
	Type   string       `json:"_type"`
	Props  BSGItemProps `json:"_props"`
}

type BSGItemProps struct {
	StackMaxSize int  `json:"StackMaxSize"`
	IsUnbuyable  bool `json:"IsUnbuyable"`
}
