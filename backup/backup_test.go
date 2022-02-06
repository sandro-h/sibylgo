package backup

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	tu "github.com/sandro-h/sibylgo/testutil"
	"github.com/sandro-h/sibylgo/util"
	"github.com/stretchr/testify/assert"
)

var origFunc func() time.Time

func TestSave(t *testing.T) {
	todoDir := tu.MakeTempDir("sibyl_backup_test")
	defer tu.DeleteTempDir(todoDir)
	todoFile := filepath.Join(todoDir, "todo.txt")
	util.WriteFile(todoFile, "my todo content 1")
	files := util.NewFileConfigFromTodoFile(todoFile)

	backup, err := Save(files, "save 1")

	assert.NoError(t, err)
	assert.Equal(t, "save 1", backup.Message)
	assert.NotNil(t, backup.Identifier)
	backups, err := ListBackups(files.TodoDir)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(backups))
	assert.Equal(t, backup, backups[0])
}

func TestSave_Multiple(t *testing.T) {
	todoDir := tu.MakeTempDir("sibyl_backup_test")
	defer tu.DeleteTempDir(todoDir)
	todoFile := filepath.Join(todoDir, "todo.txt")
	util.WriteFile(todoFile, "my todo content 1")
	files := util.NewFileConfigFromTodoFile(todoFile)

	backup1, _ := Save(files, "save 1")
	util.WriteFile(todoFile, "my todo content 2")
	backup2, err := Save(files, "save 2")

	assert.NoError(t, err)
	assert.Equal(t, "save 1", backup1.Message)
	assert.Equal(t, "save 2", backup2.Message)
	assert.NotEqual(t, backup1.Identifier, backup2.Identifier)
	backups, err := ListBackups(files.TodoDir)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(backups))
	assert.Equal(t, backup2, backups[0])
	assert.Equal(t, backup1, backups[1])
}

func TestRestore(t *testing.T) {
	// Given
	todoDir := tu.MakeTempDir("sibyl_backup_test")
	defer tu.DeleteTempDir(todoDir)
	todoFile := filepath.Join(todoDir, "todo.txt")
	util.WriteFile(todoFile, "my todo content 1")
	trashFile := filepath.Join(todoDir, "todo-trash.txt")
	util.WriteFile(trashFile, "my trash content 1")
	files := util.NewFileConfigFromTodoFile(todoFile)

	backup1, _ := Save(files, "save 1")
	util.WriteFile(files.TodoFile, "my todo content 2")
	util.WriteFile(files.TrashFile, "my trash content 2")
	Save(files, "save 2")
	util.WriteFile(files.TodoFile, "my todo content 3")
	util.WriteFile(files.TrashFile, "my trash content 3")
	Save(files, "save 3")

	// When
	restoreBackup, err := Restore(files, backup1)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Restore backup %s 'save 1'", backup1.Identifier), restoreBackup.Message)

	backups, err := ListBackups(files.TodoDir)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(backups))
	restoredTodoContent, _ := util.ReadFile(files.TodoFile)
	tu.AssertContains(t, "my todo content 1", restoredTodoContent)
	restoredTrashContent, _ := util.ReadFile(files.TrashFile)
	tu.AssertContains(t, "my trash content 1", restoredTrashContent)
}

func TestDailyBackup_NoBackupsAtAll(t *testing.T) {
	// Given
	todoDir := tu.MakeTempDir("sibyl_backup_test")
	defer tu.DeleteTempDir(todoDir)
	todoFile := filepath.Join(todoDir, "todo.txt")
	util.WriteFile(todoFile, "my todo content 1")
	files := util.NewFileConfigFromTodoFile(todoFile)
	setFakeTime("13.01.2019 12:02:42")
	defer resetOriginalTime()

	// When
	backup, err := CheckAndMakeDailyBackup(files)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "Daily backup for 13.01.2019", backup.Message)
	backups, err := ListBackups(files.TodoDir)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(backups))
	assert.Equal(t, backup, backups[0])
}

func TestDailyBackup_NoDailyBackups(t *testing.T) {
	// Given
	todoDir := tu.MakeTempDir("sibyl_backup_test")
	defer tu.DeleteTempDir(todoDir)
	todoFile := filepath.Join(todoDir, "todo.txt")
	util.WriteFile(todoFile, "my todo content 1")
	files := util.NewFileConfigFromTodoFile(todoFile)
	setFakeTime("13.01.2019 12:02:42")
	defer resetOriginalTime()
	Save(files, "some other backup")

	// When
	backup, err := CheckAndMakeDailyBackup(files)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "Daily backup for 13.01.2019", backup.Message)
	backups, err := ListBackups(files.TodoDir)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(backups))
	assert.Equal(t, backup, backups[0])
}

func TestDailyBackup_OldDailyBackups(t *testing.T) {
	// Given
	todoDir := tu.MakeTempDir("sibyl_backup_test")
	defer tu.DeleteTempDir(todoDir)
	todoFile := filepath.Join(todoDir, "todo.txt")
	util.WriteFile(todoFile, "my todo content 1")
	files := util.NewFileConfigFromTodoFile(todoFile)

	setFakeTime("12.01.2019 12:02:42")
	oldBackup, _ := CheckAndMakeDailyBackup(files)
	setFakeTime("13.01.2019 12:02:42")
	defer resetOriginalTime()

	// When
	backup, err := CheckAndMakeDailyBackup(files)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "Daily backup for 13.01.2019", backup.Message)
	backups, err := ListBackups(files.TodoDir)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(backups))
	assert.Equal(t, backup, backups[0])
	assert.Equal(t, oldBackup, backups[1])
	assert.NotEqual(t, oldBackup, backup)
}

func TestDailyBackup_AlreadyGotDailyBackup(t *testing.T) {
	// Given
	todoDir := tu.MakeTempDir("sibyl_backup_test")
	defer tu.DeleteTempDir(todoDir)
	todoFile := filepath.Join(todoDir, "todo.txt")
	util.WriteFile(todoFile, "my todo content 1")
	files := util.NewFileConfigFromTodoFile(todoFile)

	setFakeTime("13.01.2019 12:02:42")
	defer resetOriginalTime()
	oldBackup, _ := CheckAndMakeDailyBackup(files)

	// When
	backup, err := CheckAndMakeDailyBackup(files)

	// Then
	assert.NoError(t, err)
	assert.Nil(t, backup)
	backups, err := ListBackups(files.TodoDir)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(backups))
	assert.Equal(t, oldBackup, backups[0])
}

func setFakeTime(fakeTime string) {
	if origFunc == nil {
		origFunc = getNow
	}
	getNow = func() time.Time {
		t, _ := time.ParseInLocation("02.01.2006 15:04:05", fakeTime, time.Local)
		return t
	}
}

func resetOriginalTime() {
	if origFunc != nil {
		getNow = origFunc
	}
}
