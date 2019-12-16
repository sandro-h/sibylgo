package extsources

import (
	"fmt"
	tu "github.com/sandro-h/sibylgo/testutil"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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

func mockCountResponse(count int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, fmt.Sprintf(`{"count": %d}`, count))
	}))
}
