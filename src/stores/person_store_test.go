package stores

import (
	"institute-person-api/src/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPersonStore(t *testing.T) {
	cfg := config.NewConfig()
	assert.NotNil(t, cfg)
}

func TestInsert(t *testing.T) {
	cfg := config.NewConfig()
	assert.NotNil(t, cfg)
}

func TestFindOnePersonAndUpdate(t *testing.T) {
	cfg := config.NewConfig()
	assert.NotNil(t, cfg)
}
