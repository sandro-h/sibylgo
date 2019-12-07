package extsources

import (
	"fmt"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/util"
	"net/http"
	"time"
)

// FetchBitbucketPRsFromConfig calls FetchBitbucketPRs using the values in the passed config.
func FetchBitbucketPRsFromConfig(cfg *util.Config) ([]moment.Moment, error) {
	bbURL := cfg.GetString("bb_url", "")
	bbToken := cfg.GetString("bb_token", "")
	category := cfg.GetString("category", "")
	if bbURL == "" {
		return nil, fmt.Errorf("bb_url not set in config")
	}
	if bbToken == "" {
		return nil, fmt.Errorf("bb_token not set in config")
	}

	return FetchBitbucketPRs(bbURL, bbToken, category)
}

// FetchBitbucketPRs returns a single TODO moment if the user denoted by the bbToken has any open
// pull-requests in Bitbucket.
func FetchBitbucketPRs(bbBaseURL string, bbToken string, category string) ([]moment.Moment, error) {
	apiURL := fmt.Sprintf("%s/rest/api/latest/inbox/pull-requests/count", bbBaseURL)
	client := &http.Client{Timeout: 10 * time.Second}
	var count pullRequestCount
	err := util.FetchJSONAsModel(client, apiURL, &count)
	if err != nil {
		return nil, err
	}

	if count.Count == 0 {
		return nil, nil
	}

	mom := moment.NewSingleMoment(fmt.Sprintf("Reviews - %d PRs", count.Count))
	mom.SetID(&moment.Identifier{Value: "bbprs"})
	if category != "" {
		mom.SetCategory(&moment.Category{Name: category})
	}

	return []moment.Moment{mom}, nil
}

type pullRequestCount struct {
	Count int
}
