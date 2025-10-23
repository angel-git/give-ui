package api

import (
	"fmt"
	"regexp"
	"slices"
	"spt-give-ui/backend/config"
	"spt-give-ui/backend/models"
	"strings"
)

var validNameRe = regexp.MustCompile(`^[ A-Za-z0-9_-]+$`)

func CreateNewBundle(name string, description string, config *config.Config) (e error) {
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if !validNameRe.MatchString(name) {
		return fmt.Errorf("Bundle name '%s' characters which are not allowed. Use only letters, numbers or spaces", name)
	}

	currentBundles := config.GetBundles()
	bundleByNameExists := slices.ContainsFunc(currentBundles, func(bundle models.Bundle) bool {
		return bundle.Name == name
	})
	if bundleByNameExists {
		return fmt.Errorf("A Bundle already exists with the name '%s'. Please choose a different name.", name)
	}

	newBundle := models.Bundle{
		Name:        name,
		Description: description,
		Items:       map[string]uint16{},
	}
	currentBundles = append(currentBundles, newBundle)
	config.SetBundles(currentBundles)
	return nil
}
