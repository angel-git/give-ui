package models

import (
	"slices"
	"strings"
)

type Bundle struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Items       map[string]int `json:"items"`
}

type BundleItem struct {
	Id       string
	Quantity int
}

func (b *Bundle) SortedItems() []BundleItem {
	items := make([]BundleItem, 0, len(b.Items))
	for id, quantity := range b.Items {
		items = append(items, BundleItem{Id: id, Quantity: quantity})
	}

	slices.SortStableFunc(items, func(a, b BundleItem) int {
		return strings.Compare(a.Id, b.Id)
	})

	return items
}
