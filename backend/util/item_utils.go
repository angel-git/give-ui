package util

import "spt-give-ui/backend/models"

func SearchLink(item models.BSGItem) []string {
	var linkedItems []string
	// dry this
	if item.Props.Slots != nil {
		for _, slot := range *item.Props.Slots {
			if slot.Props.Filters != nil {
				for _, filter := range *slot.Props.Filters {
					if filter.Filter != nil {
						for _, f := range filter.Filter {
							linkedItems = append(linkedItems, f)
						}
					}
				}
			}
		}
	}
	if item.Props.Grids != nil {
		for _, grid := range *item.Props.Grids {
			if grid.Props.Filters != nil {
				for _, filter := range *grid.Props.Filters {
					if filter.Filter != nil {
						for _, f := range filter.Filter {
							linkedItems = append(linkedItems, f)
						}
					}
				}
			}
		}
	}
	if item.Props.Cartridges != nil {
		for _, cartridge := range *item.Props.Cartridges {
			if cartridge.Props.Filters != nil {
				for _, filter := range *cartridge.Props.Filters {
					if filter.Filter != nil {
						for _, f := range filter.Filter {
							linkedItems = append(linkedItems, f)
						}
					}
				}
			}
		}
	}
	if item.Props.Chambers != nil {
		for _, chamber := range *item.Props.Chambers {
			if chamber.Props.Filters != nil {
				for _, filter := range *chamber.Props.Filters {
					if filter.Filter != nil {
						for _, f := range filter.Filter {
							linkedItems = append(linkedItems, f)
						}
					}
				}
			}
		}
	}
	return linkedItems
}
