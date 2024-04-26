package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewBreadcrumb(t *testing.T) {
	crumb, err := NewBreadCrumb("255.255.255.255", "000000000000000000000000", "CORRID")
	assert.Nil(t, err)

	oid, err := primitive.ObjectIDFromHex("000000000000000000000000")
	assert.Nil(t, err)

	assert.NotNil(t, crumb)
	assert.Equal(t, crumb.FromIp, "255.255.255.255")
	assert.Equal(t, crumb.ByUser, oid)
	assert.Equal(t, crumb.CorrelationId, "CORRID")
	assert.WithinDuration(t, time.Now(), crumb.AtTime.Time(), time.Second)
}
