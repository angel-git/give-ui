package images

import (
	"spt-give-ui/backend/models"
	"strings"
)

type IdAndTpl struct {
	Id  string
	Tpl string
}

func GetItemSize(item models.ItemWithUpd, inventoryItems []models.ItemWithUpd, bsgItems map[string]models.BSGItem) (int, int) {
	bsgItem := bsgItems[item.Tpl]

	allChildrenItems := []IdAndTpl{}

	isItemContainer := isContainer(bsgItem)
	if !isItemContainer {
		allChildrenItems =
			findAllIdsAndTplFromParent(item.Id, inventoryItems, "mod_")
	}
	outX := bsgItem.Props.Width
	outY := bsgItem.Props.Height
	sizeUp := 0
	sizeDown := 0
	sizeLeft := 0
	sizeRight := 0
	forcedUp := 0
	forcedDown := 0
	forcedLeft := 0
	forcedRight := 0

	if bsgItem.Props.Foldable && bsgItem.Props.FoldedSlot != nil && *bsgItem.Props.FoldedSlot == "" && item.Upd.Foldable != nil && item.Upd.Foldable.Folded {
		outX -= bsgItem.Props.SizeReduceRight
	}

	for _, child := range allChildrenItems {
		childBsgItem := bsgItems[child.Tpl]
		if childBsgItem.Props.ExtraSizeForceAdd {
			forcedUp += childBsgItem.Props.ExtraSizeUp
			forcedDown += childBsgItem.Props.ExtraSizeDown
			forcedLeft += childBsgItem.Props.ExtraSizeLeft
			forcedRight += childBsgItem.Props.ExtraSizeRight
		} else {
			if sizeUp < childBsgItem.Props.ExtraSizeUp {
				sizeUp = childBsgItem.Props.ExtraSizeUp
			}
			if sizeDown < childBsgItem.Props.ExtraSizeDown {
				sizeDown = childBsgItem.Props.ExtraSizeDown
			}
			if sizeLeft < childBsgItem.Props.ExtraSizeLeft {
				sizeLeft = childBsgItem.Props.ExtraSizeLeft
			}
			if sizeRight < childBsgItem.Props.ExtraSizeRight {
				sizeRight = childBsgItem.Props.ExtraSizeRight
			}
		}

	}
	sizeX := outX + sizeLeft + sizeRight + forcedLeft + forcedRight
	sizeY := outY + sizeUp + sizeDown + forcedUp + forcedDown
	return sizeX, sizeY

}

func isContainer(bsgItem models.BSGItem) bool {
	return bsgItem.Props.Grids != nil && len(*bsgItem.Props.Grids) > 0
}

func findAllIdsAndTplFromParent(parentId string, inventoryItems []models.ItemWithUpd, slotId string) []IdAndTpl {
	var allChildrenItems []IdAndTpl
	for _, item := range inventoryItems {
		if item.ParentID != nil && *item.ParentID == parentId {
			if item.SlotID != nil && strings.HasPrefix(*item.SlotID, slotId) {
				allChildrenItems = append(allChildrenItems, IdAndTpl{Id: item.Id, Tpl: item.Tpl})
			}
			allChildrenItems = append(allChildrenItems, findAllIdsAndTplFromParent(item.Id, inventoryItems, slotId)...)
		}
	}
	return allChildrenItems
}
