package popup

import (
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"github.com/sandro-h/sibylgo/backup"
	"github.com/sandro-h/sibylgo/modify"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/util"
	log "github.com/sirupsen/logrus"
)

// Start creates the popup window and starts the event loop. The popup is not shown from the start, only
// when the hotkey is pressed.
func Start(todoFile string, popupCfg *util.Config) {
	insertCategory := popupCfg.GetString("category", "")
	if popupCfg.GetBool("dark_mode", false) {
		os.Setenv("FYNE_THEME", "dark")
	}

	a := app.New()

	w := newWindow(a)
	populateWindow(w, func(str string) {
		if str != "" {
			err := insertMoment(todoFile, insertCategory, str)
			if err != nil {
				log.Error(err)
			}
		}
		w.Hide()
	})

	go listenForHotkey(w, popupCfg.GetStringList("hotkey", nil), func(e hook.Event) {
		w.Show()
	})

	showThenHide(w)
	a.Run()
}

func populateWindow(w fyne.Window, onSubmitted func(string)) {
	entry := newTypeableEntry()
	entry.OnSubmitted = onSubmitted
	entry.onTypedKey = func(key *fyne.KeyEvent) {
		if key.Name == "Escape" {
			w.Hide()
		}
	}

	cont := container.New(layout.NewFormLayout(), widget.NewLabel("TODO"), entry)
	w.SetContent(cont)

	w.Resize(fyne.NewSize(500, cont.MinSize().Height))
	w.Canvas().Focus(entry)
	w.CenterOnScreen()
}

func newWindow(a fyne.App) fyne.Window {
	if drv, ok := a.Driver().(desktop.Driver); ok {
		return drv.CreateSplashWindow()
	}
	return a.NewWindow("")
}

func showThenHide(w fyne.Window) {
	w.Show()
	go hideImmediately(w)
}

func hideImmediately(w fyne.Window) {
	for w.Content().Visible() {
		time.Sleep(100 * time.Millisecond)
		w.Hide()
	}
}

func listenForHotkey(w fyne.Window, hotkey []string, handler func(hook.Event)) {
	if hotkey != nil {
		robotgo.EventHook(hook.KeyDown, hotkey, handler)
	}

	s := robotgo.EventStart()
	<-robotgo.EventProcess(s)
}

func insertMoment(todoFile, category string, str string) error {
	mom := moment.NewSingleMoment(str)

	if category != "" {
		mom.SetCategory(&moment.Category{Name: category})
	}

	log.Infof("Inserting '%s''\n", str)
	_, err := backup.Save(todoFile, "Backup before programmatically inserting moment")
	if err != nil {
		return err
	}

	return modify.PrependInFile(todoFile, []moment.Moment{mom})
}

type typeableEntry struct {
	widget.Entry
	onTypedKey func(key *fyne.KeyEvent)
}

func newTypeableEntry() *typeableEntry {
	e := &typeableEntry{}
	e.ExtendBaseWidget(e)
	return e
}

func (e *typeableEntry) TypedKey(key *fyne.KeyEvent) {
	e.Entry.TypedKey(key)
	if e.onTypedKey != nil {
		e.onTypedKey(key)
	}
}
