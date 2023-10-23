package tests

import (
	"institute-person-api/src/main/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBreadcrumb(t *testing.T) {
	crumb := models.NewBreadCrumb("255.255.255.255", "USERID", "CORRID")

	assert.NotNil(t, crumb)
	assert.Equal(t, crumb.FromIp, "255.255.255.255")
	assert.Equal(t, crumb.ByUser, "USERID")
	assert.Equal(t, crumb.CorrelationId, "CORRID")
	// TODO: assert time within a few seconds of now
}
