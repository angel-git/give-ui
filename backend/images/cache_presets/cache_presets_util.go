package cache_presets

import (
	"sort"
	"spt-give-ui/backend/images/cache"
	"spt-give-ui/backend/models"
)

func GetItemHash(item models.ItemWithUpd, items []models.ItemWithUpd, bsgItemsRoot map[string]models.BSGItem) int32 {
	var hash int32 = 17
	for _, h := range smethod0(item, items, bsgItemsRoot, 1) {
		hash ^= h
	}

	if cache.IsAmmoItem(item.Tpl, bsgItemsRoot) {
		hash ^= 27 * 56
	}

	return hash
}

func smethod0(topLevelItem models.ItemWithUpd, items []models.ItemWithUpd, bsgItemsRoot map[string]models.BSGItem, hashSeed int32) []int32 {
	var hashes []int32

	hashes = append(hashes, smethod1(topLevelItem, items, bsgItemsRoot)*hashSeed)

	if cache.IsHideEntrails(topLevelItem.Tpl, bsgItemsRoot) {
		return hashes
	}

	hashSeed *= 6529
	childrenItems := getChildren(topLevelItem, items)
	if len(childrenItems) > 0 {
		for _, child := range childrenItems {
			var num int32 = 0
			var num2 = hashSeed ^ getHashSum(child, items)
			if isSlackSlot(topLevelItem.Tpl, bsgItemsRoot) {
				num++
				num2 ^= 2879 * num
			}
			hashes = append(hashes, smethod0(child, items, bsgItemsRoot, num2)...)
		}
	}

	sort.Slice(hashes, func(i, j int) bool {
		return hashes[i] < hashes[j]
	})
	return hashes
}

func getHashSum(item models.ItemWithUpd, items []models.ItemWithUpd) int32 {
	parentItem := getParentItem(item, items)
	containerID := item.SlotID
	num := 2777 * cache.GetDeterministicHashCode(*containerID)
	num += 7901 * cache.GetHashCodeFromMongoID(parentItem.Tpl)
	return num
}

func isSlackSlot(tpl string, bsgItemsRoot map[string]models.BSGItem) bool {
	container := bsgItemsRoot[tpl]
	if container.Props.Cartridges != nil {
		return len(*container.Props.Cartridges) > 0
	} else {
		return false
	}
}

func smethod1(item models.ItemWithUpd, items []models.ItemWithUpd, bsgItemsRoot map[string]models.BSGItem) int32 {
	hash := cache.GetHashCodeFromMongoID(item.Tpl)

	node := bsgItemsRoot[item.Tpl]
	if node.Props.HasHinge {
		isToggled := true
		if item.Upd != nil && item.Upd.Togglable != nil {
			isToggled = item.Upd.Togglable.On
		}
		hash ^= 23 + cache.BoolToInt(isToggled)
	}

	if cache.IsFoldableItem(item.Tpl, bsgItemsRoot) {
		isFolded := false
		if item.Upd != nil && item.Upd.Foldable != nil {
			isFolded = item.Upd.Foldable.Folded
		}
		hash ^= (23 + cache.BoolToInt(isFolded)) << 1
	}

	if cache.IsMagazineItem(item.Tpl, bsgItemsRoot) {
		allAmmoInsideMagazine := 0
		for _, i := range items {
			if i.ParentID != nil && *i.ParentID == item.Id {
				if i.Upd != nil {
					allAmmoInsideMagazine += i.Upd.StackObjectsCount
				} else {
					allAmmoInsideMagazine++
				}
			}
		}

		maxVisibleAmmo := cache.GetMaxVisibleAmmo(uint16(allAmmoInsideMagazine), node.Props.VisibleAmmoRangesString)
		hash ^= (23 + int32(maxVisibleAmmo)) << 2
	}

	return hash
}

func getParentItem(item models.ItemWithUpd, items []models.ItemWithUpd) *models.ItemWithUpd {
	for _, i := range items {
		if i.Id == *item.ParentID {
			return &i
		}
	}
	return nil
}

func getChildren(item models.ItemWithUpd, items []models.ItemWithUpd) []models.ItemWithUpd {
	var children []models.ItemWithUpd
	for _, i := range items {
		if i.ParentID != nil && *i.ParentID == item.Id {
			children = append(children, i)
		}
	}
	return children
}
