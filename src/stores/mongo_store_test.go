package stores

import (
	"institute-person-api/src/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMongoStore(t *testing.T) {
	cfg := config.NewConfig()
	assert.NotNil(t, cfg)
}

func TestInsertOne(t *testing.T) {
	cfg := config.NewConfig()
	assert.NotNil(t, cfg)
}

func TestFindOne(t *testing.T) {
	cfg := config.NewConfig()
	assert.NotNil(t, cfg)
}

func TestFindMany(t *testing.T) {
	cfg := config.NewConfig()
	assert.NotNil(t, cfg)
}

func TestFindOneAndUpdate(t *testing.T) {
	cfg := config.NewConfig()
	assert.NotNil(t, cfg)
}

func TestFindId(t *testing.T) {
	cfg := config.NewConfig()
	assert.NotNil(t, cfg)
}

func TestFindDocuments(t *testing.T) {
	cfg := config.NewConfig()
	assert.NotNil(t, cfg)
}
