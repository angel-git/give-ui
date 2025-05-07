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
	StackMaxSize            int                `json:"StackMaxSize"`
	IsUnbuyable             bool               `json:"IsUnbuyable"`
	HasHinge                bool               `json:"HasHinge"`
	Foldable                bool               `json:"Foldable"`
	FoldedSlot              *string            `json:"FoldedSlot"`
	VisibleAmmoRangesString string             `json:"VisibleAmmoRangesString"`
	HideEntrails            bool               `json:"HideEntrails"`
	Cartridges              *[]WithFilterProps `json:"Cartridges"`
	Chambers                *[]WithFilterProps `json:"Chambers"`
	Grids                   *[]WithFilterProps `json:"Grids"`
	Width                   int                `json:"Width"`
	Height                  int                `json:"Height"`
	SizeReduceRight         int                `json:"SizeReduceRight"`
	ExtraSizeForceAdd       bool               `json:"ExtraSizeForceAdd"`
	ExtraSizeUp             int                `json:"ExtraSizeUp"`
	ExtraSizeDown           int                `json:"ExtraSizeDown"`
	ExtraSizeLeft           int                `json:"ExtraSizeLeft"`
	ExtraSizeRight          int                `json:"ExtraSizeRight"`
	BackgroundColor         string             `json:"BackgroundColor"`
	Slots                   *[]WithFilterProps `json:"Slots"`
}

type WithFilterProps struct {
	Props FilterProps `json:"_props"`
}

type FilterProps struct {
	Filters *[]Filters `json:"filters"`
}

type Filters struct {
	Filter []string `json:"Filter"`
}
