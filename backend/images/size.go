package images

import (
	"spt-give-ui/backend/models"
	"strings"
)

func GetItemSize(rootItem models.ItemWithUpd, inventoryItems []models.ItemWithUpd, bsgItems map[string]models.BSGItem) (int, int) {
	bsgItem := bsgItems[rootItem.Tpl]

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

	bsgItemFoldedSlot := ""
	if bsgItem.Props.FoldedSlot != nil {
		bsgItemFoldedSlot = *bsgItem.Props.FoldedSlot
	}
	isBsgItemFoldable := bsgItem.Props.Foldable
	isRootFolded := rootItem.Upd != nil && rootItem.Upd.Foldable != nil && rootItem.Upd.Foldable.Folded
	if isBsgItemFoldable && (bsgItem.Props.FoldedSlot == nil || *bsgItem.Props.FoldedSlot == "") && isRootFolded {
		outX -= bsgItem.Props.SizeReduceRight
	}

	toDo := []string{rootItem.Id}
	for len(toDo) > 0 {
		id := toDo[0]
		toDo = toDo[1:]

		if isContainer(bsgItem) {
			continue
		}

		allChildrenItems := findAllIdsAndTplFromParent(id, inventoryItems)
		for _, child := range allChildrenItems {
			if child.SlotID == nil || !strings.HasPrefix(*child.SlotID, "mod_") {
				continue
			}

			toDo = append(toDo, child.Id)

			childBsgItem := bsgItems[child.Tpl]
			childFoldable := childBsgItem.Props.Foldable
			childFolded := child.Upd != nil && child.Upd.Foldable != nil && child.Upd.Foldable.Folded

			if isBsgItemFoldable && bsgItemFoldedSlot == *child.SlotID && (isRootFolded || childFolded) {
				continue
			}

			if childFoldable && isRootFolded && childFolded {
				continue
			}

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

	}
	sizeX := outX + sizeLeft + sizeRight + forcedLeft + forcedRight
	sizeY := outY + sizeUp + sizeDown + forcedUp + forcedDown
	return sizeX, sizeY

}

func isContainer(bsgItem models.BSGItem) bool {
	return bsgItem.Props.Grids != nil && len(*bsgItem.Props.Grids) > 0
}

func findAllIdsAndTplFromParent(parentId string, inventoryItems []models.ItemWithUpd) []models.ItemWithUpd {
	var allChildrenItems []models.ItemWithUpd
	for _, item := range inventoryItems {
		if item.ParentID != nil && *item.ParentID == parentId {
			allChildrenItems = append(allChildrenItems, item)
		}
	}
	return allChildrenItems
}
