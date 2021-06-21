package mode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModes(t *testing.T) {
	assert.Truef(t, IsDebug(), "mode should be %s by default", Debug)

	SetMode(Prod)

	assert.Truef(t, IsProd(), "mode should be %s but is %s", Prod, CurrentMode())

	SetMode(-1) //non-existent

	assert.Equal(t, "Unknown", CurrentMode().String())
}
