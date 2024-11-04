# Fit - Git Wrapper CLI

**Fit** is a minimalistic Git wrapper written in Go that enhances the readability of Git command outputs. It includes formatting improvements, color-coded outputs, and enforces English output for consistency across different system languages.

## Features

- Executes common Git commands like `status`, `log`, `branch`, and `diff` with a clearer, formatted output.
- Enforces Git output in English to avoid inconsistencies across different language settings.
- Highlights errors, warnings, and status messages in a structured, color-coded format.

## Installation

1. Clone the repository.
2. Build the executable:
   ```bash
   go build -o fit
   ```

## Usage

Run the `fit` command followed by any Git command arguments:

```bash
./fit status     # Displays formatted git status
./fit log        # Shows git log in a structured table
./fit branch     # Lists branches with the current branch highlighted
./fit diff       # Shows diff with color-coded changes
```

## Example Outputs

- **Status**: Shows file modifications in a table with labels for modified, new, and deleted files.
- **Log**: Lists commits with date, hash, author, and message in an easy-to-read table.
- **Branch**: Displays branches, with a marker for the active branch.
- **Diff**: Color-codes added, removed, and modified lines for clarity.

## Notes

Fit sets the `LANG` environment variable to English (`en_US.UTF-8`), ensuring consistent output in English regardless of the system's language setting.
