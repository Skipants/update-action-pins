package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/otiai10/copy"
)

func TestCorrectFile(t *testing.T) {
	tmpDir := "tmp"
	err := copy.Copy("test/fixtures", tmpDir)
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	for _, filePair := range [][]string{
		{"tmp/workflows/good_workflow.yml", "tmp/workflows/good_workflow.yml"},
		{"tmp/workflows/no_deps_workflow.yml", "tmp/workflows/no_deps_workflow.yml"},
		{"tmp/workflows/bad_workflow.yml", "tmp/workflows/bad_workflow_fixed.yml"},
	} {
		actualFilename, expectedFilename := filePair[0], filePair[1]

		err = correctFile(actualFilename)

		actual, err := os.ReadFile(actualFilename)
		if err != nil {
			t.Fatalf("failed to read original file: %v", err)
		}
		expected, err := os.ReadFile(expectedFilename)
		if err != nil {
			t.Fatalf("failed to read expected file: %v", err)
		}

		if !bytes.Equal(actual, expected) {
			t.Fatalf("actual file %s does not match expected file %s", actualFilename, expectedFilename)
		}
	}
}
