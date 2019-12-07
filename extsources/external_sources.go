package extsources

import (
	"fmt"
	"github.com/sandro-h/sibylgo/util"
	"time"
)

// ExternalSourcesProcess periodically checks a list of external sources (based on the passed config) for moments.
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
	fmt.Printf("%s\n", p.extSrcConfig.GetSubConfig("bitbucket_prs").GetString("bb_url", ""))
}
