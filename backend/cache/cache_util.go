package cache

import (
	"spt-give-ui/backend/models"
	"strconv"
	"strings"
)

func GetItemHash(item models.InventoryItem, items []models.InventoryItem, bsgItemsRoot map[string]interface{}) int32 {
	var hash int32 = 17
	for _, h := range smethod0(item, items, bsgItemsRoot, 1) {
		hash ^= h
	}

	if isAmmoItem(item.Tpl, bsgItemsRoot) {
		hash ^= 27 * 56
	}

	return hash
}

func smethod0(topLevelItem models.InventoryItem, items []models.InventoryItem, bsgItemsRoot map[string]interface{}, hashSeed int32) []int32 {
	var hashes []int32

	hashes = append(hashes, smethod1(topLevelItem, items, bsgItemsRoot)*hashSeed)

	if isHideEntrails(topLevelItem.Tpl, bsgItemsRoot) {
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

	return hashes
}

func getHashSum(item models.InventoryItem, items []models.InventoryItem) int32 {
	parentItem := getParentItem(item, items)
	containerID := item.SlotID
	num := 2777 * getDeterministicHashCode(*containerID)
	num += 7901 * getHashCodeFromMongoID(parentItem.Tpl)
	return num
}

func isSlackSlot(tpl string, bsgItemsRoot map[string]interface{}) bool {
	container := bsgItemsRoot[tpl].(map[string]interface{})
	if cartridges, ok := container["_props"].(map[string]interface{})["Cartridges"].([]interface{}); ok {
		return len(cartridges) > 0
	}
	return false
}

func smethod1(item models.InventoryItem, items []models.InventoryItem, bsgItemsRoot map[string]interface{}) int32 {
	hash := getHashCodeFromMongoID(item.Tpl)

	node := bsgItemsRoot[item.Tpl].(map[string]interface{})
	if hasHinge, ok := node["_props"].(map[string]interface{})["HasHinge"].(bool); ok && hasHinge {
		isToggled := true
		if item.Upd != nil && item.Upd.Togglable != nil {
			isToggled = item.Upd.Togglable.On
		}
		hash ^= 23 + boolToInt(isToggled)
	}

	if isFoldableItem(item.Tpl, bsgItemsRoot) {
		isFolded := false
		if item.Upd != nil && item.Upd.Foldable != nil {
			isFolded = item.Upd.Foldable.Folded
		}
		hash ^= (23 + boolToInt(isFolded)) << 1
	}

	if isMagazineItem(item.Tpl, bsgItemsRoot) {
		allAmmoInsideMagazine := 0
		for _, i := range items {
			if i.ParentID != nil && *i.ParentID == item.ID {
				if i.Upd != nil {
					allAmmoInsideMagazine += i.Upd.StackObjectsCount
				} else {
					allAmmoInsideMagazine++
				}
			}
		}

		maxVisibleAmmo := getMaxVisibleAmmo(uint16(allAmmoInsideMagazine), node["_props"].(map[string]interface{})["VisibleAmmoRangesString"].(string))
		hash ^= (23 + int32(maxVisibleAmmo)) << 2
	}

	return hash
}

func getMaxVisibleAmmo(bullets uint16, visibleAmmoRangesString string) uint16 {
	visibleAmmoRanges := getMaxVisibleAmmoRanges(visibleAmmoRangesString)

	for i, r := range visibleAmmoRanges {
		if r[0] <= bullets && r[1] >= bullets {
			return bullets
		}
		if bullets < r[0] {
			if i == 0 {
				return r[0]
			}
			return visibleAmmoRanges[i-1][1]
		}
	}

	return visibleAmmoRanges[len(visibleAmmoRanges)-1][1]
}

func getMaxVisibleAmmoRanges(visibleAmmoRangesString string) [][2]uint16 {
	if visibleAmmoRangesString == "" {
		return [][2]uint16{{1, 2}}
	}

	var ranges [][2]uint16
	splits := strings.Split(visibleAmmoRangesString, ";")
	for _, split := range splits {
		rangeParts := strings.Split(split, "-")
		start, _ := strconv.ParseUint(rangeParts[0], 10, 16)
		end, _ := strconv.ParseUint(rangeParts[1], 10, 16)
		ranges = append(ranges, [2]uint16{uint16(start), uint16(end)})
	}

	return ranges
}

func getDeterministicHashCode(s string) int32 {
	var hash1 int32 = 5381
	var hash2 int32 = hash1

	for i := 0; i < len(s); i += 2 {
		hash1 = (hash1 << 5) + hash1 ^ int32(s[i])

		if i+1 < len(s) {
			hash2 = (hash2 << 5) + hash2 ^ int32(s[i+1])
		}
	}

	return hash1 + hash2*1566083941
}

func getHashCodeFromMongoID(mongoID string) int32 {
	timestamp := getMongoIDTimestamp(mongoID)
	counter := getMongoIDCounter(mongoID)
	counterHigh := int32((counter >> 32) * 3637)
	counterLow := int32((counter & 0xFFFFFFFF) * 5807)

	return int32(timestamp) ^ counterHigh ^ counterLow
}

func getMongoIDTimestamp(mongoID string) uint32 {
	timestamp, err := strconv.ParseUint(mongoID[:8], 16, 32)
	if err != nil {
		panic(err)
	}
	return uint32(timestamp)
}

func getMongoIDCounter(mongoID string) uint64 {
	counter, err := strconv.ParseUint(mongoID[8:24], 16, 64)
	if err != nil {
		panic(err)
	}
	return counter
}

func isAmmoItem(tpl string, bsgItemsRoot map[string]interface{}) bool {
	return findParentByName(bsgItemsRoot, tpl, "Ammo") != nil
}

func isMagazineItem(tpl string, bsgItemsRoot map[string]interface{}) bool {
	return findParentByName(bsgItemsRoot, tpl, "Magazine") != nil
}

func isFoldableItem(tpl string, bsgItemsRoot map[string]interface{}) bool {
	node := bsgItemsRoot[tpl].(map[string]interface{})
	if foldable, ok := node["_props"].(map[string]interface{})["Foldable"].(bool); ok {
		return foldable
	}
	return false
}

func getParentItem(item models.InventoryItem, items []models.InventoryItem) models.InventoryItem {
	for _, i := range items {
		if i.ID == *item.ParentID {
			return i
		}
	}
	return nil
}

func isHideEntrails(tpl string, bsgItemsRoot map[string]interface{}) bool {
	node := bsgItemsRoot[tpl].(map[string]interface{})
	if hideEntrails, ok := node["_props"].(map[string]interface{})["HideEntrails"].(bool); ok {
		return hideEntrails
	}
	return false
}

func getChildren(item models.InventoryItem, items []models.InventoryItem) []models.InventoryItem {
	var children []models.InventoryItem
	for _, i := range items {
		if i.ParentID != nil && *i.ParentID == item.ID {
			children = append(children, i)
		}
	}
	return children
}

func findParentByName(bsgItemsRoot map[string]interface{}, currentID, targetName string) interface{} {
	node := bsgItemsRoot[currentID].(map[string]interface{})
	name := node["_name"].(string)
	if name == targetName {
		return node
	}
	if parentID, ok := node["_parent"].(string); ok {
		return findParentByName(bsgItemsRoot, parentID, targetName)
	}
	return nil
}

func boolToInt(b bool) int32 {
	if b {
		return 1
	}
	return 0
}
