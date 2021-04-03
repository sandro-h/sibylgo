package testutil

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var updateGolden = flag.Bool("update-golden", false, "Update golden test files")
var dryGolden = flag.Bool("dry-golden", false, "Used together with -update-golden. If set, write to separate, ignored files to compare manually.")

// AssertGoldenOutput reads the expected output of a test from goldenFile and checks that
// the actual ouput matches.
func AssertGoldenOutput(t *testing.T, testName string, goldenFile string, output string) {
	if *updateGolden {
		outputFile := goldenFile
		if *dryGolden {
			outputFile += ".tmp"
		}
		WriteTestdata(t, testName, outputFile, output)
	} else {
		expectedOutput := ReadTestdata(t, testName, goldenFile)
		assert.Equal(t, expectedOutput, output, "Testcase: %s", testName)
	}
}

// ReadTestdata reads a file from the current package's testdata/ folder.
func ReadTestdata(t *testing.T, testName string, path string) string {
	data, err := readFile(FullTestdataPath(path))
	if err != nil {
		assert.Failf(t, "Could not load testdata", "Testcase: %s, file: %s, error: %s", testName, path, err)
	}
	return strings.ReplaceAll(data, "\r", "")
}

// WriteTestdata wrotes a file to the current package's testdata/ folder.
func WriteTestdata(t *testing.T, testName string, path string, data string) {
	err := writeFile(FullTestdataPath(path), data)
	if err != nil {
		assert.Failf(t, "Could not write testdata", "Testcase: %s, file: %s, error: %s", testName, path, err)
	}
}

// FullTestdataPath returns the full path of a testdata file, given a partial path (usually just the filename).
// E.g.: GetTestdataFile('my_testdata.txt') -> testdata/my_testdata.txt
func FullTestdataPath(filenameOrPartialPath string) string {
	return filepath.Join("testdata", filenameOrPartialPath)
}

func readFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func writeFile(filePath string, str string) error {
	if str[len(str)-1] != '\n' {
		str += "\n"
	}
	return os.WriteFile(filePath, []byte(str), 0644)
}
