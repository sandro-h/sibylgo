package outlook

import (
	"errors"
	"fmt"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/parse"
	"github.com/sandro-h/sibylgo/util"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

// CheckInfinitely repeatedly checks for changes in the todo file and updates
// Outlook events.
func CheckInfinitely(todoFile string, interval time.Duration) {
	lastMod := time.Unix(0, 0)
	for {
		CheckOnce(todoFile, &lastMod)
		time.Sleep(interval)
	}
}

// CheckOnce checks for changes in the todo file and updates
// Outlook events.
func CheckOnce(todoFile string, lastMod *time.Time) {
	file, err := os.Stat(todoFile)
	if err != nil {
		log.Errorf("Could not stat %s: %s\n", todoFile, err)
		return
	}

	newLastMod := file.ModTime()
	if newLastMod == *lastMod {
		return
	}

	*lastMod = newLastMod
	todos, err := parse.File(todoFile)
	if err != nil {
		log.Errorf("Could not read todo file %s: %s\n", todoFile, err)
		return
	}

	err = UpdateOutlookEvents(todos.Moments)
	if err != nil {
		log.Errorf("Had one or more errors updating outlook events: %s\n", err)
	}
}

// UpdateOutlookEvents syncs single moments with specific date to Outlook as events.
func UpdateOutlookEvents(moments []moment.Moment) error {
	currentMoms := filterEligibleForOutlook(moments)
	outlookMoms, err := listEvents()
	if err != nil {
		return fmt.Errorf("Could not list outlook events: %s", err)
	}

	debugMomentList("Current moments", currentMoms)
	debugMomentList("Outlook moments", outlookMoms)

	addedMoments, updatedMoments, removedMoments := computeDiff(currentMoms, outlookMoms)
	debugMomentList("Added moments", addedMoments)
	debugMomentList("Updated moments", updatedMoments)
	debugMomentList("Removed moments", removedMoments)

	allErrors := ""
	err = removeEvents(removedMoments, updatedMoments)
	if err != nil {
		allErrors += err.Error() + "\n"
	}

	err = createEvents(addedMoments, updatedMoments)
	if err != nil {
		allErrors += err.Error() + "\n"
	}

	if allErrors != "" {
		return errors.New(allErrors)
	}

	return nil
}

func filterEligibleForOutlook(moments []moment.Moment) []*moment.SingleMoment {
	var res []*moment.SingleMoment
	for _, m := range moments {
		// Only single, non-range moments
		singMom, ok := m.(*moment.SingleMoment)
		if ok && !singMom.IsDone() {

			if moment.IsSingleDayMoment(singMom) {
				res = append(res, singMom)
			} else if moment.IsDueMoment(singMom) {
				// Due moment must be converted to single day moment, since we don't encode "due" in outlook.
				// (we don't want long ranged outlook events)
				clone := *singMom
				clone.Start = &moment.Date{Time: util.SetToStartOfDay(singMom.End.Time)}
				clone.End = &moment.Date{Time: util.SetToEndOfDay(singMom.End.Time)}
				res = append(res, &clone)
			}
		}
	}
	return res
}

func computeDiff(currentMoms []*moment.SingleMoment, outlookMoms []*moment.SingleMoment) ([]*moment.SingleMoment, []*moment.SingleMoment, []*moment.SingleMoment) {
	currentMomsByName := groupByName(currentMoms)
	outlookMomsByName := groupByName(outlookMoms)

	var addedMoments []*moment.SingleMoment
	var updatedMoments []*moment.SingleMoment
	var removedMoments []*moment.SingleMoment

	for _, m := range currentMoms {
		outlookMom, found := outlookMomsByName[m.GetName()]
		if found {
			if !momEqual(m, outlookMom) {
				updatedMoments = append(updatedMoments, m)
			}
		} else {
			addedMoments = append(addedMoments, m)
		}
	}

	for _, m := range outlookMoms {
		_, found := currentMomsByName[m.GetName()]
		if !found {
			removedMoments = append(removedMoments, m)
		}
	}

	return addedMoments, updatedMoments, removedMoments
}

func groupByName(moms []*moment.SingleMoment) map[string]*moment.SingleMoment {
	res := make(map[string]*moment.SingleMoment)
	for _, m := range moms {
		res[m.GetName()] = m
	}
	return res
}

func momEqual(a *moment.SingleMoment, b *moment.SingleMoment) bool {
	if a.GetName() != b.GetName() {
		return false
	}

	if !a.Start.Time.Equal(b.Start.Time) {
		return false
	}

	aToD := time.Unix(0, 0)
	bToD := time.Unix(0, 0)
	if a.TimeOfDay != nil {
		aToD = a.TimeOfDay.Time
	}
	if b.TimeOfDay != nil {
		bToD = b.TimeOfDay.Time
	}
	if !aToD.Equal(bToD) {
		return false
	}

	return true
}

func debugMomentList(header string, list []*moment.SingleMoment) {
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("%s:\n", header)
		for _, m := range list {
			log.Debugf("  %s\n", momentToDebugString(m))
		}
	}
}

func momentToDebugString(a *moment.SingleMoment) string {
	aToD := "-"
	if a.TimeOfDay != nil {
		aToD = a.TimeOfDay.Time.String()
	}

	return fmt.Sprintf("name=%-32s start=%-32s timeOfDay=%-32s", a.GetName(), a.Start.Time, aToD)
}

func removeEvents(listOfLists ...[]*moment.SingleMoment) error {
	allErrors := ""
	for _, list := range listOfLists {
		for _, m := range list {
			err := removeEvent(m)
			if err != nil {
				allErrors += fmt.Sprintf("Could not remove %s: %s\n", m.GetName(), err)
			}
		}
	}

	if allErrors != "" {
		return errors.New(allErrors)
	}

	return nil
}

func createEvents(listOfLists ...[]*moment.SingleMoment) error {
	allErrors := ""
	for _, list := range listOfLists {
		for _, m := range list {
			err := createEvent(m)
			if err != nil {
				allErrors += fmt.Sprintf("Could not create %s: %s\n", m.GetName(), err)
			}
		}
	}

	if allErrors != "" {
		return errors.New(allErrors)
	}

	return nil
}
