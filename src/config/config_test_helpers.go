package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function to assert configuration item properties
func assertConfigItem(t *testing.T, cfg *Config, itemName, expectedFrom string, expectedValue interface{}) {

	FindConfigItem := func(items []*ConfigItem, name string) int {
		for index, item := range items {
			if item.Name == name {
				return index
			}
		}
		return -1 // Return -1 when the item is not found
	}

	itemIndex := FindConfigItem(cfg.ConfigItems, itemName)
	assert.NotEqual(t, -1, itemIndex, "Config item not found: "+itemName)
	if itemIndex != -1 {
		assert.Equal(t, itemName, cfg.ConfigItems[itemIndex].Name)
		assert.Equal(t, expectedFrom, cfg.ConfigItems[itemIndex].From)
		assert.Equal(t, expectedValue, cfg.ConfigItems[itemIndex].Value)
	}
}
