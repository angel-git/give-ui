package commands

import (
	"fmt"
	"spt-give-ui/backend/models"
)

func AddItem(itemId string, amount int) models.Command {
	return models.Command{
		Message: fmt.Sprintf("spt give %s %d", itemId, amount),
	}
}

func AddUserPreset(itemId string) models.Command {
	return models.Command{
		Message: fmt.Sprintf("spt give-user-preset %s", itemId),
	}
}

func AddStashItem(itemId string) models.Command {
	return models.Command{
		Message: fmt.Sprintf("spt give-user-stash-item %s", itemId),
	}
}

func AddGearPreset(presetId string) models.Command {
	return models.Command{
		Message: fmt.Sprintf("spt give-gear-preset %s", presetId),
	}
}

func UpdateTraderRep(nickname string, rep string) models.Command {
	return models.Command{
		Message: fmt.Sprintf("spt trader %s rep %s", nickname, rep),
	}
}

func UpdateTraderSpend(nickname string, spend string) models.Command {
	return models.Command{
		Message: fmt.Sprintf("spt trader %s spend %s", nickname, spend),
	}
}

func UpdateLevel(level int) models.Command {
	return models.Command{
		Message: fmt.Sprintf("spt profile level %d", level),
	}
}

func UpdateSkill(skill string, progress int) models.Command {
	return models.Command{
		Message: fmt.Sprintf("spt profile skill %s %d", skill, progress),
	}
}

func SetWinterSeason() models.Command {
	return models.Command{
		Message: "itsonlysnowalan",
	}
}

func SetSummerSeason() models.Command {
	return models.Command{
		Message: "givemesunshine",
	}
}

func SetHalloweenSeason() models.Command {
	return models.Command{
		Message: "veryspooky",
	}
}

func SetChristmasSeason() models.Command {
	return models.Command{
		Message: "hohoho",
	}
}

func AddRowsToStash() models.Command {
	return models.Command{
		Message: "givemespace",
	}
}

func Gift(gifId string) models.Command {
	return models.Command{
		Message: gifId,
	}
}
