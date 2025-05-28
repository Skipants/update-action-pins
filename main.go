package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var helpFlags = []string{"--help", "-h"}

func main() {
	if len(os.Args) < 2 || slices.Contains(helpFlags, os.Args[1]) {
		fmt.Println("Usage: update-action-pins <file-or-dir>")
		return
	}

	fileOrDirPath := os.Args[1]
	fileOrDirInfo, err := os.Stat(os.Args[1])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var files []string

	if fileOrDirInfo.IsDir() {
		err = filepath.Walk(fileOrDirPath, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !fi.IsDir() {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	} else {
		files = append(files, fileOrDirPath)
	}

	var shaFromActionVersion = func(action string, version string) (string, error) {
		// todo replace this with github api stuff
		return "", nil
	}

	for _, file := range files {
		if err := correctFile(file, shaFromActionVersion); err != nil {
			fmt.Println("Error processing", file, ":", err)
		}
	}
}

func correctFile(filename string, shaFromActionVersion func(string, string) (string, error)) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	usesRegex := regexp.MustCompile(`uses:\s*([^\s@]+)@([^\s]+)`)
	shaRegex := regexp.MustCompile(`^[0-9a-fA-F]{40}$`)
	for scanner.Scan() {
		currLine := scanner.Text()
		matches := usesRegex.FindStringSubmatch(currLine)

		if matches != nil {
			action := matches[1]
			version := matches[2]

			if !shaRegex.MatchString(version) {
				sha, err := shaFromActionVersion(action, version)
				if err != nil {
					return fmt.Errorf("couldn't get a sha for the line: %s: %w", strings.TrimSpace(currLine), err)
				}

				currLine = usesRegex.ReplaceAllString(currLine, fmt.Sprintf("uses: %s@%s", action, sha))
			}
		}
		lines = append(lines, currLine)
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	file.Close()
	file, err = os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, line := range lines {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
