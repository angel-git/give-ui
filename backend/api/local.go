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
		Items:       map[string]int{},
	}
	currentBundles = append(currentBundles, newBundle)
	config.SetBundles(currentBundles)
	return nil
}

func AddItemToBundle(bundleName string, itemID string, quantity int, config *config.Config) (e error) {
	bundles := config.GetBundles()
	bundleIndex := slices.IndexFunc(bundles, func(bundle models.Bundle) bool {
		return bundle.Name == bundleName
	})
	if bundleIndex == -1 {
		return fmt.Errorf("No Bundle found with the name '%s'", bundleName)
	}

	bundle := &bundles[bundleIndex]
	if existingQty, exists := bundle.Items[itemID]; exists {
		bundle.Items[itemID] = existingQty + quantity
	} else {
		bundle.Items[itemID] = quantity
	}

	config.SetBundles(bundles)
	return nil
}

func RemoveItemFromBundle(bundleName string, itemID string, config *config.Config) (e error) {
	bundles := config.GetBundles()
	bundleIndex := slices.IndexFunc(bundles, func(bundle models.Bundle) bool {
		return bundle.Name == bundleName
	})
	if bundleIndex == -1 {
		return fmt.Errorf("No Bundle found with the name '%s'", bundleName)
	}

	bundle := &bundles[bundleIndex]
	if _, exists := bundle.Items[itemID]; !exists {
		return fmt.Errorf("Item with ID '%s' not found in Bundle '%s'", itemID, bundleName)
	}

	delete(bundle.Items, itemID)
	config.SetBundles(bundles)
	return nil
}

func DeleteBundle(bundleName string, config *config.Config) (e error) {
	bundles := config.GetBundles()
	bundleIndex := slices.IndexFunc(bundles, func(bundle models.Bundle) bool {
		return bundle.Name == bundleName
	})
	if bundleIndex == -1 {
		return fmt.Errorf("No Bundle found with the name '%s'", bundleName)
	}

	bundles = append(bundles[:bundleIndex], bundles[bundleIndex+1:]...)
	config.SetBundles(bundles)
	return nil
}
