package netbox

import (
	"testing"

	"github.com/netbox-community/go-netbox/netbox/models"
	"github.com/stretchr/testify/assert"
)

func TestGetTagListFromNestedTagList(t *testing.T) {

	tags := []*models.NestedTag{
		&models.NestedTag{
			Name: strToPtr("Foo"),
			Slug: strToPtr("foo"),
		},
		&models.NestedTag{
			Name: strToPtr("Bar"),
			Slug: strToPtr("bar"),
		},
	}

	flat := getTagListFromNestedTagList(tags)
	expected := []string{
		"Foo",
		"Bar",
	}
	assert.Equal(t, flat, expected)
}
