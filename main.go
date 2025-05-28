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
	if len(os.Args) < 2 {
		fmt.Println("Usage: update-action-pins <file-or-dir>")
		return
	}
	root := os.Args[1]
	var files []string
	info, err := os.Stat(root)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if info.IsDir() {
		err = filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
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
		files = append(files, root)
	}
	for _, file := range files {
		if err := correctFile(file); err != nil {
			fmt.Println("Error processing", file, ":", err)
		}
	}
}

func correctFile(filename string) error {
	return nil
	file, err := os.Open(filename)
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
				fmt.Printf("%s:%d: %s\n", filename, lineNum, strings.TrimSpace(line))
			}
		}
		lineNum++
	}
	return scanner.Err()
}
