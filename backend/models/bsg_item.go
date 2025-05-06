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
	StackMaxSize            int     `json:"StackMaxSize"`
	IsUnbuyable             bool    `json:"IsUnbuyable"`
	HasHinge                bool    `json:"HasHinge"`
	Foldable                bool    `json:"Foldable"`
	FoldedSlot              *string `json:"FoldedSlot"`
	VisibleAmmoRangesString string  `json:"VisibleAmmoRangesString"`
	HideEntrails            bool    `json:"HideEntrails"`
	Cartridges              *[]any  `json:"Cartridges"`
	Grids                   *[]Grid `json:"Grids"`
	Width                   int     `json:"Width"`
	Height                  int     `json:"Height"`
	SizeReduceRight         int     `json:"SizeReduceRight"`
	ExtraSizeForceAdd       bool    `json:"ExtraSizeForceAdd"`
	ExtraSizeUp             int     `json:"ExtraSizeUp"`
	ExtraSizeDown           int     `json:"ExtraSizeDown"`
	ExtraSizeLeft           int     `json:"ExtraSizeLeft"`
	ExtraSizeRight          int     `json:"ExtraSizeRight"`
	BackgroundColor         string  `json:"BackgroundColor"`
	Slots                   *[]Slot `json:"Slots"`
}

type Slot struct {
	Props SlotProps `json:"_props"`
}

type Grid struct {
	Props GridProps `json:"_props"`
}

type GridProps struct {
	Filters *[]Filters `json:"filters"`
}

type SlotProps struct {
	Filters *[]Filters `json:"filters"`
}

type Filters struct {
	Filter []string `json:"Filter"`
}
