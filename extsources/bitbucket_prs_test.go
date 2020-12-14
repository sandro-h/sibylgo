package extsources

import (
	tu "github.com/sandro-h/sibylgo/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFetchBitbucketPRs(t *testing.T) {
	ts := tu.MockSimpleJSONResponse(`{"count": 12}`)
	defer ts.Close()

	moms, err := FetchBitbucketPRs(ts.URL, "myuser", "1234", "Today")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(moms))
	assert.Equal(t, "Reviews - 12 PRs", moms[0].GetName())
	assert.Equal(t, "bbprs", moms[0].GetID().Value)
	assert.Equal(t, "Today", moms[0].GetCategory().Name)
}

func TestFetchBitbucketPRs_NoPRs(t *testing.T) {
	ts := tu.MockSimpleJSONResponse(`{"count": 0}`)
	defer ts.Close()

	moms, err := FetchBitbucketPRs(ts.URL, "myuser", "1234", "Today")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(moms))
}
