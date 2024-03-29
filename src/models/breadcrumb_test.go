package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBreadcrumb(t *testing.T) {
	crumb, err := NewBreadCrumb("255.255.255.255", "000000000000000000000000", "CORRID")

	assert.NotNil(t, err)
	assert.NotNil(t, crumb)
	assert.Equal(t, crumb.FromIp, "255.255.255.255")
	assert.Equal(t, crumb.ByUser, "USERID")
	assert.Equal(t, crumb.CorrelationId, "CORRID")
	// TODO: assert time within a few seconds of now
}
