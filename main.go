package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/go-github/v72/github"
	"github.com/gregjones/httpcache"
	"github.com/urfave/cli/v3"
)

var CURRENT_VERSION = "0.0.0"

func main() {
	cmd := &cli.Command{
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return updateActionPins(cmd.StringArg("fileOrDirPath"))
		},
		UsageText: `update-action-pins [global options] <file-or-directory-path>

<file-or-directory-path> is the path to the github action files you would like to run this against. It defaults to ".github/workflows" if no argument is given.`,
		Arguments: []cli.Argument{
			&cli.StringArg{
				Config: cli.StringConfig{
					TrimSpace: true,
				},
				Name:  "fileOrDirPath",
				Value: ".github/workflows",
			},
		},
		Name:    "update-action-pins",
		Version: CURRENT_VERSION,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func updateActionPins(fileOrDirPath string) error {
	fileOrDirInfo, err := os.Stat(fileOrDirPath)
	if err != nil {
		return err
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
			return err
		}
	} else {
		files = append(files, fileOrDirPath)
	}

	githubClient := github.NewClient(httpcache.NewMemoryCacheTransport().Client()).WithAuthToken(os.Getenv("GITHUB_TOKEN"))

	var shaFromActionVersion = func(action string, version string) (string, error) {
		parts := strings.Split(action, "/")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid action format: %s", action)
		}
		owner, repo := parts[0], parts[1]

		ref, _, err := githubClient.Git.GetRef(context.Background(), owner, repo, "refs/heads/"+version)
		if err == nil && ref.Object != nil {
			return ref.Object.GetSHA(), nil
		}

		ref, _, err = githubClient.Git.GetRef(context.Background(), owner, repo, "refs/tags/"+version)
		if err == nil && ref.Object != nil {
			sha := ref.Object.GetSHA()
			if ref.Object.GetType() == "tag" {
				tagObj, _, tagErr := githubClient.Git.GetTag(context.Background(), owner, repo, sha)
				if tagErr == nil && tagObj.Object != nil {
					return tagObj.Object.GetSHA(), nil
				}
			}
			return sha, nil
		}

		return "", fmt.Errorf("could not find branch or tag '%s' for %s/%s", version, owner, repo)
	}

	for _, file := range files {
		if err := correctFile(file, shaFromActionVersion); err != nil {
			fmt.Println("Error processing", file, ":", err)
		}
	}

	return nil
}

func correctFile(filename string, shaFromActionVersion func(string, string) (string, error)) error {
	if !strings.HasSuffix(filename, ".yml") && !strings.HasSuffix(filename, ".yaml") {
		return nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	isWorkflow := false
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "on:") || strings.HasPrefix(trimmed, "jobs:") {
			isWorkflow = true
			break
		}
	}
	if !isWorkflow {
		return nil
	}

	file.Seek(0, 0)

	var lines []string
	scanner = bufio.NewScanner(file)
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
					fmt.Println("Warning:", fmt.Errorf("couldn't get a sha for the line: %s: %w", strings.TrimSpace(currLine), err))
				} else {
					currLine = usesRegex.ReplaceAllString(currLine, fmt.Sprintf("uses: %s@%s # %s", action, sha, version))
				}
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
