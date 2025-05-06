package util

import "spt-give-ui/backend/models"

func SearchLink(item models.BSGItem) []string {
	var linkedItems []string

	//for _, v := range allItems {
	if item.Props.Slots != nil {
		for _, slot := range *item.Props.Slots {
			if slot.Props.Filters != nil {
				for _, filter := range *slot.Props.Filters {
					if filter.Filter != nil {
						for _, f := range filter.Filter {
							//return append(SearchLink(v, allItems), v)
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
	//}
	return linkedItems
}
