package extsources

import (
	"fmt"
	"github.com/sandro-h/sibylgo/backup"
	"github.com/sandro-h/sibylgo/modify"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/parse"
	"github.com/sandro-h/sibylgo/util"
	"strings"
	"time"
)

var externalSources = map[string]fetchFunc{
	"dummies":       FetchDummyMomentsFromConfig,
	"bitbucket_prs": FetchBitbucketPRsFromConfig,
}

const idPrefix = "ext_"

type fetchFunc func(*util.Config) ([]moment.Moment, error)

// ExternalSourcesProcess periodically checks a list of external sources (based on the passed config) for moments,
// then updates the todo file with them.
type ExternalSourcesProcess struct {
	todoFilePath  string
	extSrcConfig  *util.Config
	checkInterval time.Duration
}

// NewExternalSourcesProcess creates a new ExternalSourcesProcess.
func NewExternalSourcesProcess(todoFilePath string, extSrcConfig *util.Config) *ExternalSourcesProcess {
	return &ExternalSourcesProcess{todoFilePath, extSrcConfig, 10 * time.Minute}
}

// CheckInfinitely repeatedly checks the external sources in the check interval.
// This method blocks indefinitely and should be run as a go routine.
func (p *ExternalSourcesProcess) CheckInfinitely() {
	for {
		p.CheckOnce()
		time.Sleep(p.checkInterval)
	}
}

// CheckOnce does a single check on the external sources.
func (p *ExternalSourcesProcess) CheckOnce() {
	content, err := util.ReadFile(p.todoFilePath)
	if err != nil {
		fmt.Printf("[Ext sources] Failed to read todo file %s: %s\n", p.todoFilePath, err.Error())
	}

	updatedContent, err := FetchAndApplyExternalSourceMoments(content, p.extSrcConfig)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	if updatedContent != content {
		backup.Save(p.todoFilePath, "Backup before applying external source changes")
		err = util.WriteFile(p.todoFilePath, updatedContent)
		if err != nil {
			fmt.Printf("[Ext sources] Failed to write todo file %s: %s\n", p.todoFilePath, err.Error())
		}
	}
}

// FetchAndApplyExternalSourceMoments fetches all moments from the configured sources and
// updates the passed todo file content with them. It adds or updates (based on ID) all
// moments found in the external sources. It removes any moments with an ext_* ID
// that were not found in any external source anymore.
func FetchAndApplyExternalSourceMoments(content string, extSrcConfig *util.Config) (string, error) {
	fetchedMoments, fetchedMomentsByID := fetchExternalSourceMoments(extSrcConfig)

	todos, err := parse.String(content)
	if err != nil {
		return "", fmt.Errorf("[Ext sources] Failed to parse todo file: %s", err.Error())
	}

	var toDelete []moment.Moment
	for id, m := range todos.MomentsByID {
		if strings.HasPrefix(id, idPrefix) {
			_, found := fetchedMomentsByID[id]
			if !found {
				toDelete = append(toDelete, m)
			}
		}
	}

	// TODO could be optimized in one modify call that does removes and upserts.
	prepend := extSrcConfig.GetBool("prepend", false)
	content, _ = modify.Delete(content, toDelete)
	content, err = modify.Upsert(content, fetchedMoments, prepend)
	if err != nil {
		return "", fmt.Errorf("[Ext sources] Failed to upsert moments: %s", err.Error())
	}
	return content, nil
}

func fetchExternalSourceMoments(extSrcConfig *util.Config) ([]moment.Moment, map[string]moment.Moment) {
	var allMoments []moment.Moment
	byID := make(map[string]moment.Moment)
	for srcName, fetchFunc := range externalSources {
		if extSrcConfig.HasKey(srcName) {
			moments, err := fetchFunc(extSrcConfig.GetSubConfig(srcName))
			if err != nil {
				fmt.Printf("[Ext sources] Fetching %s failed: %s\n", srcName, err.Error())
			} else {
				for _, m := range moments {
					if m.GetID() == nil {
						fmt.Printf("[Ext sources] %s moment %s has no ID. Skipping it.\n", srcName, m.GetName())
						continue
					}

					fullID := idPrefix + m.GetID().Value
					m.GetID().Value = fullID
					allMoments = append(allMoments, m)
					byID[fullID] = m
				}
			}
		}
	}
	return allMoments, byID
}
