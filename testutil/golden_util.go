package testutil

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
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
	data, err := readFile(filepath.Join("testdata", path))
	if err != nil {
		assert.Failf(t, "Could not load testdata", "Testcase: %s, file: %s, error: %s", testName, path, err)
	}
	return strings.ReplaceAll(data, "\r", "")
}

// WriteTestdata wrotes a file to the current package's testdata/ folder.
func WriteTestdata(t *testing.T, testName string, path string, data string) {
	err := writeFile(filepath.Join("testdata", path), data)
	if err != nil {
		assert.Failf(t, "Could not write testdata", "Testcase: %s, file: %s, error: %s", testName, path, err)
	}
}

func readFile(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func writeFile(filePath string, str string) error {
	if str[len(str)-1] != '\n' {
		str += "\n"
	}
	return ioutil.WriteFile(filePath, []byte(str), 0644)
}
