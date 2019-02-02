package util

import (
	tu "github.com/sandro-h/sibylgo/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetToStartOfWeek(t *testing.T) {
	assert.Equal(t, "28.01.2019 00:00:00", tu.Dtts(SetToStartOfWeek(tu.Dt("02.02.2019"))))
	// Sunday
	assert.Equal(t, "28.01.2019 00:00:00", tu.Dtts(SetToStartOfWeek(tu.Dt("03.02.2019"))))
	// Monday
	assert.Equal(t, "04.02.2019 00:00:00", tu.Dtts(SetToStartOfWeek(tu.Dt("04.02.2019"))))
}

func TestSetToEndOfWeek(t *testing.T) {
	assert.Equal(t, "03.02.2019 23:59:59", tu.Dtts(SetToEndOfWeek(tu.Dt("02.02.2019"))))
	// Sunday
	assert.Equal(t, "03.02.2019 23:59:59", tu.Dtts(SetToEndOfWeek(tu.Dt("03.02.2019"))))
	// Monday
	assert.Equal(t, "10.02.2019 23:59:59", tu.Dtts(SetToEndOfWeek(tu.Dt("04.02.2019"))))
}
