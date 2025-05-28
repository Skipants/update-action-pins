package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	err := filepath.Walk(".github/workflows", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !(strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".yaml")) {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		re := regexp.MustCompile(`^\s*uses:\s*([^\s@]+)@([^\s]+)`)
		shaRe := regexp.MustCompile(`^[0-9a-fA-F]{40}$`)
		lineNum := 1
		for scanner.Scan() {
			line := scanner.Text()
			matches := re.FindStringSubmatch(line)
			if len(matches) == 3 {
				version := matches[2]
				if !shaRe.MatchString(version) {
					fmt.Printf("%s:%d: %s\n", path, lineNum, strings.TrimSpace(line))
				}
			}
			lineNum++
		}
		return scanner.Err()
	})

	if err != nil {
		fmt.Println("Error:", err)
	}
}
