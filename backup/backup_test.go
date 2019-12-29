package backup

import (
	"fmt"
	tu "github.com/sandro-h/sibylgo/testutil"
	"github.com/sandro-h/sibylgo/util"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestSave(t *testing.T) {
	todoDir := tu.MakeTempDir("sibyl_backup_test")
	defer tu.DeleteTempDir(todoDir)
	todoFile := filepath.Join(todoDir, "todo.txt")
	util.WriteFile(todoFile, "my todo content 1")

	backup, err := Save(todoFile, "save 1")

	assert.NoError(t, err)
	assert.Equal(t, "save 1", backup.Message)
	assert.NotNil(t, backup.Identifier)
	backups, err := ListBackups(todoFile)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(backups))
	assert.Equal(t, backup, backups[0])
}

func TestSave_Multiple(t *testing.T) {
	todoDir := tu.MakeTempDir("sibyl_backup_test")
	defer tu.DeleteTempDir(todoDir)
	todoFile := filepath.Join(todoDir, "todo.txt")
	util.WriteFile(todoFile, "my todo content 1")

	backup1, err := Save(todoFile, "save 1")
	util.WriteFile(todoFile, "my todo content 2")
	backup2, err := Save(todoFile, "save 2")

	assert.NoError(t, err)
	assert.Equal(t, "save 1", backup1.Message)
	assert.Equal(t, "save 2", backup2.Message)
	assert.NotEqual(t, backup1.Identifier, backup2.Identifier)
	backups, err := ListBackups(todoFile)
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

	backup1, err := Save(todoFile, "save 1")
	util.WriteFile(todoFile, "my todo content 2")
	Save(todoFile, "save 2")
	util.WriteFile(todoFile, "my todo content 3")
	Save(todoFile, "save 3")

	// When
	restoreBackup, err := Restore(todoFile, backup1)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Restore backup %s 'save 1'", backup1.Identifier), restoreBackup.Message)

	backups, err := ListBackups(todoFile)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(backups))
	restoredContent, _ := util.ReadFile(todoFile)
	tu.AssertContains(t, "my todo content 1", restoredContent)
}
