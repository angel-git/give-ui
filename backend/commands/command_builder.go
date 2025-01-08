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

func UpdateTraderRep(nickname string, rep string) models.Command {
	return models.Command{
		Message: fmt.Sprintf("spt trader %s %s", nickname, rep),
	}
}

func UpdateTraderSpend(nickname string, spend string) models.Command {
	return models.Command{
		Message: fmt.Sprintf("spt trader %s %s", nickname, spend),
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
