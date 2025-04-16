package templ_utils

import (
	"fmt"
	"spt-give-ui/backend/models"
	"strings"
)

type InventoryItem struct {
	Item  models.ItemWithUpd
	PosX  int
	PosY  int
	SizeX int
	SizeY int
}

func PrepareItemsForStash(profile *models.SPTProfile) ([]*InventoryItem, int, int) {
	inventoryItems := []*InventoryItem{}
	maxX, maxY := 0, 0

	allItems := profile.Characters.PMC.Inventory.Items
	stashId := profile.Characters.PMC.Inventory.Stash
	for i := range allItems {
		item := &allItems[i]
		if item.ParentID != nil && *item.ParentID == stashId && item.Location != nil {
			inventoryItem := &InventoryItem{
				Item:  *item,
				PosX:  item.Location.X,
				PosY:  item.Location.Y,
				SizeX: item.SizeX,
				SizeY: item.SizeY,
			}
			if item.Location.R == "Vertical" {
				inventoryItem.SizeY = item.SizeX
				inventoryItem.SizeX = item.SizeY
			}
			inventoryItems = append(inventoryItems, inventoryItem)
			if inventoryItem.PosX+inventoryItem.SizeX > maxX {
				maxX = inventoryItem.PosX + inventoryItem.SizeX
			}
			if inventoryItem.PosY+inventoryItem.SizeY > maxY {
				maxY = inventoryItem.PosY + inventoryItem.SizeY
			}
		}
	}

	gridSizeX := 10
	if maxX > gridSizeX {
		gridSizeX = maxX
	}
	gridSizeY := 50
	if maxY > gridSizeY {
		gridSizeY = maxY
	}

	grid := make([][]*InventoryItem, gridSizeY)
	for i := range grid {
		grid[i] = make([]*InventoryItem, gridSizeX) // Initialize with nil
	}

	for _, item := range inventoryItems {
		locX := item.PosX
		locY := item.PosY
		sizeX := item.SizeX
		sizeY := item.SizeY
		for col := locY; col < locY+sizeY; col++ {
			for row := locX; row < locX+sizeX; row++ {
				// Boundary check
				if col >= 0 && col < gridSizeY && row >= 0 && row < gridSizeX {
					if grid[col][row] != nil {
						// Overlap detected! Handle error? Log? Ignore?
						// For now, let's just log it. Depending on Tarkov logic,
						// this might indicate bad data.
						fmt.Printf("Warning: Overlap detected at [%d][%d] for item %s\n", col, row, item.Item.Id)
					}
					grid[col][row] = item
				} else {
					fmt.Printf("Warning: Item %s (%dx%d at %d,%d) partially out of bounds (%dx%d grid)\n",
						item.Item.Id, sizeX, sizeY, locX, locY, gridSizeX, gridSizeY)
				}
			}
		}
	}

	// 4. Flatten the grid into the ordered list
	orderedItems := make([]*InventoryItem, 0, gridSizeX*gridSizeY)
	addedItems := make(map[string]struct{}) // Set to track added items

	for col := 0; col < gridSizeY; col++ {
		for row := 0; row < gridSizeX; row++ {
			item := grid[col][row]
			if item != nil {
				if _, found := addedItems[item.Item.Id]; !found {
					orderedItems = append(orderedItems, item) // Add the actual item
					addedItems[item.Item.Id] = struct{}{}     // Mark as added
				} else {
					orderedItems = append(orderedItems, nil) // Add placeholder for spanned cell
				}
			} else {
				orderedItems = append(orderedItems, nil) // Add placeholder for empty cell
			}
		}
	}

	return orderedItems, gridSizeX, gridSizeY
}

func buildBackgroundStyleForImage(item InventoryItem) map[string]string {
	rotation := "0deg"
	if item.Item.Location.R == "Vertical" {
		rotation = "90deg"
	}
	width := item.SizeX * 64
	if item.Item.Location.R == "Vertical" {
		width = item.SizeY * 64
	}
	height := item.SizeY * 64
	if item.Item.Location.R == "Vertical" {
		height = item.SizeX * 64
	}
	translateX := (64 - width) / 2
	translateY := (64 - height) / 2

	translateYAdjustment := (height - 64) / 2
	if rotation == "90deg" {
		translateYAdjustment = (width - 64) / 2
	}
	translateXAdjustment := (width - 64) / 2
	if rotation == "90deg" {
		translateXAdjustment = (height - 64) / 2
	}

	// Update the translateY value.

	translateYFinal := translateY + translateYAdjustment
	translateXFinal := translateX + translateXAdjustment
	styles := map[string]string{
		"background-repeat": "no-repeat",
		"width":             fmt.Sprintf("%dpx", width),
		"height":            fmt.Sprintf("%dpx", height),
		"transform":         fmt.Sprintf("rotate(%s)", rotation),
	}

	if rotation == "90deg" {
		styles["translate"] = fmt.Sprintf("%dpx %dpx", translateXFinal, translateYFinal)
	}

	if item.Item.ImageBase64 != "" {
		// add background image
		styles["background-image"] = fmt.Sprintf("url(%s)", item.Item.ImageBase64)
	} else {
		backgroundImageUrl := fmt.Sprintf("https://assets.tarkov.dev/%s-base-image.webp", item.Item.Tpl)
		styles["background-image"] = fmt.Sprintf("url(%s)", backgroundImageUrl)
	}
	return styles
}

func CalculateBackgroundStyleForImage(item InventoryItem) string {
	styles := buildBackgroundStyleForImage(item)

	var styleParts []string
	for key, value := range styles {
		styleParts = append(styleParts, fmt.Sprintf("%s: %s", key, value))
	}

	cssString := strings.Join(styleParts, "; ")
	return cssString
}

func CalculateBackgroundStyleForItem(item InventoryItem) string {
	return fmt.Sprintf("z-index: 2; position: relative; height: %dpx; width: %dpx", item.SizeY*64, item.SizeX*64)
}

//
//export const calculateBackgroundColor = (backgroundColor: string) => {
//switch (backgroundColor) {
//case 'black':
//return `rgba(0, 0, 0, 0.3)`;
//case 'blue':
//return `rgba(28, 65, 86, 0.3)`;
//case 'green':
//return `rgba(21, 45, 0, 0.3)`;
//case 'grey':
//return `rgba(29, 29, 29, 0.3)`;
//case 'orange':
//return `rgba(60, 25, 0, 0.3)`;
//case 'red':
//return `rgba(109, 36, 24, 0.3)`;
//case 'violet':
//return `rgba(76, 42, 85, 0.3)`;
//case 'yellow':
//return `rgba(104, 102, 40, 0.3)`;
//default:
//return `rgba(127, 127, 127, 0.0)`;
//}
//};
