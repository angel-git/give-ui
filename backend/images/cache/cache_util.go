package cache

import (
	"iter"
	"maps"
	"spt-give-ui/backend/models"
	"strconv"
	"strings"
)

func GetItemHash(item models.BSGItem, bsgItemsRoot map[string]models.BSGItem) int32 {
	var hash int32 = 17
	for _, h := range smethod0(item, bsgItemsRoot, 1) {
		hash ^= h
	}

	if IsAmmoItem(item.Id, bsgItemsRoot) {
		hash ^= 27 * 56
	}

	return hash
}

func smethod0(topLevelItem models.BSGItem, bsgItemsRoot map[string]models.BSGItem, hashSeed int32) []int32 {
	var hashes []int32

	hashes = append(hashes, smethod1(topLevelItem, bsgItemsRoot)*hashSeed)

	if IsHideEntrails(topLevelItem.Id, bsgItemsRoot) {
		return hashes
	}

	bsgItemsRootValues := maps.Values(bsgItemsRoot)

	hashSeed *= 6529
	childrenItems := getChildren(topLevelItem, bsgItemsRootValues)
	if len(childrenItems) > 0 {
		for _, child := range childrenItems {
			var num int32 = 0
			var num2 = hashSeed ^ getHashSum(child, bsgItemsRootValues)
			if isSlackSlot(topLevelItem.Id, bsgItemsRoot) {
				num++
				num2 ^= 2879 * num
			}
			hashes = append(hashes, smethod0(child, bsgItemsRoot, num2)...)
		}
	}

	return hashes
}

func getHashSum(item models.BSGItem, items iter.Seq[models.BSGItem]) int32 {
	parentItem := getParentItem(item, items)
	containerID := "hideout"
	num := 2777 * GetDeterministicHashCode(containerID)
	num += 7901 * GetHashCodeFromMongoID(parentItem.Id)
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

func smethod1(item models.BSGItem, bsgItemsRoot map[string]models.BSGItem) int32 {
	hash := GetHashCodeFromMongoID(item.Id)

	node := bsgItemsRoot[item.Id]
	if node.Props.HasHinge {
		isToggled := true
		//if item.Upd != nil && item.Upd.Togglable != nil {
		//	isToggled = item.Upd.Togglable.On
		//}
		hash ^= 23 + BoolToInt(isToggled)
	}

	if IsFoldableItem(item.Id, bsgItemsRoot) {
		isFolded := false
		//if item.Upd != nil && item.Upd.Foldable != nil {
		//	isFolded = item.Upd.Foldable.Folded
		//}
		hash ^= (23 + BoolToInt(isFolded)) << 1
	}

	if IsMagazineItem(item.Id, bsgItemsRoot) {
		allAmmoInsideMagazine := 0
		//for _, i := range items {
		//	if i.ParentID != nil && *i.ParentID == item.ID {
		//		if i.Upd != nil {
		//			allAmmoInsideMagazine += i.Upd.StackObjectsCount
		//		} else {
		//			allAmmoInsideMagazine++
		//		}
		//	}
		//}

		maxVisibleAmmo := GetMaxVisibleAmmo(uint16(allAmmoInsideMagazine), node.Props.VisibleAmmoRangesString)
		hash ^= (23 + int32(maxVisibleAmmo)) << 2
	}

	return hash
}

func GetMaxVisibleAmmo(bullets uint16, visibleAmmoRangesString string) uint16 {
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

func GetDeterministicHashCode(s string) int32 {
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

func GetHashCodeFromMongoID(mongoID string) int32 {
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

func IsAmmoItem(tpl string, bsgItemsRoot map[string]models.BSGItem) bool {
	return findParentByName(bsgItemsRoot, tpl, "Ammo") != nil
}

func IsMagazineItem(tpl string, bsgItemsRoot map[string]models.BSGItem) bool {
	return findParentByName(bsgItemsRoot, tpl, "Magazine") != nil
}

func IsFoldableItem(tpl string, bsgItemsRoot map[string]models.BSGItem) bool {
	node := bsgItemsRoot[tpl]
	return node.Props.Foldable
}

func getParentItem(item models.BSGItem, items iter.Seq[models.BSGItem]) *models.BSGItem {
	for i := range items {
		if i.Id == item.Parent {
			return &i
		}
	}
	return nil
}

func IsHideEntrails(tpl string, bsgItemsRoot map[string]models.BSGItem) bool {
	node := bsgItemsRoot[tpl]
	return node.Props.HideEntrails
}

func getChildren(item models.BSGItem, items iter.Seq[models.BSGItem]) []models.BSGItem {
	var children []models.BSGItem
	for i := range items {
		if i.Parent != "" && i.Parent == item.Id {
			children = append(children, i)
		}
	}
	return children
}

func findParentByName(bsgItemsRoot map[string]models.BSGItem, currentID string, targetName string) *models.BSGItem {
	node := bsgItemsRoot[currentID]
	name := node.Name
	if name == targetName {
		return &node
	}
	parent := node.Parent
	if parent == "" {
		return nil
	}
	return findParentByName(bsgItemsRoot, node.Parent, targetName)
}

func BoolToInt(b bool) int32 {
	if b {
		return 1
	}
	return 0
}
