package models

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
