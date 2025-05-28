package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/otiai10/copy"
	"gopkg.in/yaml.v3"
)

type ActionVersions struct {
	Actions map[string]map[string]string `yaml:"actions"`
}

var mockedShaFromActionVersion = func(action string, version string) (string, error) {
	file, err := os.ReadFile("test/fixtures/action-version-mocks.yml")
	if err != nil {
		return "", err
	}

	var data ActionVersions
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		return "", err
	}

	if versions, ok := data.Actions[action]; ok {
		if sha, ok := versions[version]; ok {
			return sha, nil
		}
	}
	return "", os.ErrNotExist
}

func TestCorrectFile(t *testing.T) {
	tmpDir := "tmp"
	err := copy.Copy("test/fixtures", tmpDir)
	if err != nil {
		panic(err)
	}

	for _, filePair := range [][]string{
		{"tmp/workflows/good_workflow.yml", "tmp/workflows/good_workflow.yml"},
		{"tmp/workflows/no_deps_workflow.yml", "tmp/workflows/no_deps_workflow.yml"},
		{"tmp/workflows/bad_workflow.yml", "tmp/workflows/bad_workflow_fixed.yml"},
	} {
		actualFilename, expectedFilename := filePair[0], filePair[1]

		err = correctFile(actualFilename, mockedShaFromActionVersion)
		if err != nil {
			t.Errorf("correctFile errored: %v", err)
		}

		actual, err := os.ReadFile(actualFilename)
		if err != nil {
			t.Fatalf("failed to read original file: %v", err)
		}
		expected, err := os.ReadFile(expectedFilename)
		if err != nil {
			t.Fatalf("failed to read expected file: %v", err)
		}

		if !bytes.Equal(actual, expected) {
			t.Errorf("actual file %s does not match expected file %s", actualFilename, expectedFilename)
		}
	}
}
