package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRoleNames(t *testing.T) {
	u := &User{}

	roles := u.RoleNames()
	assert.Empty(t, roles)

	u.Roles = []Role{
		{Name: "admin"},
		{Name: "engineering"},
	}

	roles = u.RoleNames()
	assert.NotEmpty(t, roles)
	assert.Len(t, roles, 2)
	assert.Equal(t, roles[0], "admin")
	assert.Equal(t, roles[1], "engineering")

}
