# update-action-pins

A command-line tool to automatically update GitHub Actions workflow files, replacing version tags or unpinned actions with pinned commit SHA identifiers. This helps ensure your workflows are reproducible and secure by preventing unexpected changes in third-party actions.

## Installation

1. Build the binary:
   ```sh
   go build -o update-action-pins main.go
   ```
2. (Optional) Move it to your PATH:
   ```sh
   sudo mv update-action-pins /usr/local/bin/
   ```

## Usage

```sh
update-action-pins <file-or-dir>
```

- `<file-or-dir>`: Path to a workflow YAML file or a directory containing workflow files. Defaults to ".github/workflows"

Example:
```sh
update-action-pins

update-action-pins .github/workflows/test.yml
```

## Requirements
- Go 1.20+
- A valid GitHub token in the `GITHUB_TOKEN` environment variable (for API requests)

## How It Works
- The tool parses each workflow file and looks for `uses:` steps.
- For each action using a version tag or branch, it queries the GitHub API to resolve the corresponding commit SHA.
- The workflow file is updated in-place, replacing the version with the resolved SHA and adding a comment with the original version.

## Example
Before:
```yaml
- uses: actions/checkout@v3
- uses: actions/setup-node@main
```
After running the tool:
```yaml
- uses: actions/checkout@b4ffde3b8c7e7e3b6b7e3e1e3b6b7e3e1e3b6b7e # v3
- uses: actions/setup-node@c4c1b6b5e2e3b6b7e3e1e3b6b7e3e1e3b6b7e3e1 # main
```

## Testing
Run the test suite with:
```sh
go test
```

## Contributing
Pull requests and issues are welcome!

## License
MIT License
