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

func TestEpochWeek(t *testing.T) {
	assert.Equal(t, 2600, EpochWeek(tu.DtUtc("31.10.2019")))
	assert.Equal(t, 2600, EpochWeek(tu.DtUtc("01.11.2019")))
	assert.Equal(t, 2600, EpochWeek(tu.DtUtc("02.11.2019")))
	assert.Equal(t, 2600, EpochWeek(tu.DtUtc("03.11.2019")))
	assert.Equal(t, 2600, EpochWeek(tu.DtUtc("04.11.2019")))
	assert.Equal(t, 2600, EpochWeek(tu.DtUtc("05.11.2019")))
	assert.Equal(t, 2600, EpochWeek(tu.DtUtc("06.11.2019")))
	assert.Equal(t, 2601, EpochWeek(tu.DtUtc("07.11.2019")))
	assert.Equal(t, 2601, EpochWeek(tu.DtUtc("08.11.2019")))
	assert.Equal(t, 2602, EpochWeek(tu.DtUtc("15.11.2019")))

	// Year boundary
	assert.Equal(t, 2608, EpochWeek(tu.DtUtc("28.12.2019")))
	assert.Equal(t, 2609, EpochWeek(tu.DtUtc("04.01.2020")))
}
