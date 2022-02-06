package extsources

import (
	"fmt"
	"strings"
	"time"

	"github.com/sandro-h/sibylgo/backup"
	"github.com/sandro-h/sibylgo/modify"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/parse"
	"github.com/sandro-h/sibylgo/util"
	log "github.com/sirupsen/logrus"
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
	files         *util.FileConfig
	extSrcConfig  *util.Config
	checkInterval time.Duration
}

// NewExternalSourcesProcess creates a new ExternalSourcesProcess.
func NewExternalSourcesProcess(files *util.FileConfig, extSrcConfig *util.Config) *ExternalSourcesProcess {
	return &ExternalSourcesProcess{files, extSrcConfig, 10 * time.Minute}
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
	content, err := util.ReadFile(p.files.TodoFile)
	if err != nil {
		log.Errorf("[Ext sources] Failed to read todo file %s: %s\n", p.files.TodoFile, err.Error())
	}

	updatedContent, err := FetchAndApplyExternalSourceMoments(content, p.extSrcConfig)
	if err != nil {
		log.Errorf("%s\n", err.Error())
	}

	if updatedContent != content {
		// Avoid backup noise because of missing trailing newlines. But we still want to write
		// the newlines to the todo file.
		if !util.EqualsIgnoreTrailingNewlines(updatedContent, content) {
			backup.Save(p.files, "Backup before applying external source changes")
		}

		err = util.WriteFile(p.files.TodoFile, updatedContent)
		if err != nil {
			log.Errorf("[Ext sources] Failed to write todo file %s: %s\n", p.files.TodoFile, err.Error())
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
		if !extSrcConfig.HasKey(srcName) {
			continue
		}
		moments, err := fetchFunc(extSrcConfig.GetSubConfig(srcName))
		if err != nil {
			log.Errorf("[Ext sources] Fetching %s failed: %s\n", srcName, err.Error())
			continue
		}
		for _, m := range moments {
			if m.GetID() == nil {
				log.Errorf("[Ext sources] %s moment %s has no ID. Skipping it.\n", srcName, m.GetName())
				continue
			}

			fullID := idPrefix + m.GetID().Value
			m.GetID().Value = fullID
			allMoments = append(allMoments, m)
			byID[fullID] = m
		}
	}
	return allMoments, byID
}
