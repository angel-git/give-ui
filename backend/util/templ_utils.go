package util

import (
	"fmt"
	"regexp"
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
	width := item.SizeX * 64
	height := item.SizeY * 64

	if item.Item.Location.R == "Vertical" {
		rotation = "90deg"
		width = item.SizeY * 64
		height = item.SizeX * 64
	}
	translateX := (64 - width) / 2
	translateY := (64 - height) / 2

	translateYAdjustment := (height - 64) / 2
	translateXAdjustment := (width - 64) / 2

	if rotation == "90deg" {
		translateYAdjustment = (width - 64) / 2
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
		styles["background-image"] = fmt.Sprintf("url(data:image/png;base64,%s)", item.Item.ImageBase64)
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
	return fmt.Sprintf("z-index: 2; position: relative; height: %dpx; width: %dpx; background-color: %s", item.SizeY*64, item.SizeX*64, calculateBackgroundColor(item))
}

func SumAllItems(templateId string, profile *models.SPTProfile) string {
	total := 0
	for _, item := range profile.Characters.PMC.Inventory.Items {
		if item.Tpl == templateId && item.Upd != nil {
			total += item.Upd.StackObjectsCount
		}
	}
	return formatWithSpaces(total)
}

func formatWithSpaces(n int) string {
	s := fmt.Sprintf("%d", n)
	var result []byte
	count := 0

	for i := len(s) - 1; i >= 0; i-- {
		result = append([]byte{s[i]}, result...)
		count++
		if count%3 == 0 && i != 0 {
			result = append([]byte{' '}, result...)
		}
	}
	return string(result)
}

func calculateBackgroundColor(item InventoryItem) string {
	color := "rgba(127, 127, 127, 0.0)"
	switch item.Item.BackgroundColor {
	case "black":
		color = "rgba(0, 0, 0, 0.3)"
	case "blue":
		color = "rgba(28, 65, 86, 0.3)"
	case "green":
		color = "rgba(21, 45, 0, 0.3)"
	case "grey":
		color = "rgba(29, 29, 29, 0.3)"
	case "orange":
		color = "rgba(60, 25, 0, 0.3)"
	case "red":
		color = "rgba(109, 36, 24, 0.3)"
	case "violet":
		color = "rgba(76, 42, 85, 0.3)"
	case "yellow":
		color = "rgba(104, 102, 40, 0.3)"
	}
	return color
}

// https://hub.sp-tarkov.com/files/file/2841-odt-s-item-info-3-11-update-added-colored-name
func ApplyOdtColors(input string) string {
	re := regexp.MustCompile(`(?i)<b><color=([#a-z0-9]+)>(.*?)</color></b>(.*)`)
	output := re.ReplaceAllString(input, `<div style="color: $1">$2 $3</div>`)
	return output
}

func RemoveOdtColors(input string) string {
	re := regexp.MustCompile(`(?i)<b><color=([#a-z0-9]+)>(.*?)</color></b>(.*)`)
	output := re.ReplaceAllString(input, `$2`)
	return output
}
