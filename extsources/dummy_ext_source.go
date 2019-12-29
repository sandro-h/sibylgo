package extsources

import (
	"fmt"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/util"
	"strings"
)

// FetchDummyMomentsFromConfig returns dummy ext source moments based on the string list
// "dummy_moments" in the config. Each dummy moment must have the format <id>:<name>.
func FetchDummyMomentsFromConfig(cfg *util.Config) ([]moment.Moment, error) {
	category := cfg.GetString("category", "")
	dummies := cfg.GetStringList("dummy_moments", nil)
	fmt.Printf("%s\n", dummies)
	var moments []moment.Moment
	for _, d := range dummies {
		parts := strings.Split(d, ":")
		id := parts[0]
		name := parts[1]
		mom := moment.NewSingleMoment(name)
		mom.SetID(&moment.Identifier{Value: id})
		if category != "" {
			mom.SetCategory(&moment.Category{Name: category})
		}
		moments = append(moments, mom)
	}
	return moments, nil
}
