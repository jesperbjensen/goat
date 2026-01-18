# Installation Guide

## Quick Install

From within the goat project directory:

```bash
go install
```

This installs the `goat` binary to `$GOPATH/bin` (typically `~/go/bin`).

## Verify Installation

Check that goat was installed:

```bash
ls $(go env GOPATH)/bin/goat
```

You should see the goat binary listed.

## Add to PATH

If you can't run `goat` from anywhere, you need to add Go's bin directory to your PATH.

### Check if already in PATH

```bash
echo $PATH | grep -q "$(go env GOPATH)/bin" && echo "Already in PATH" || echo "Not in PATH"
```

### Add to PATH (if needed)

#### For Bash

Add to `~/.bashrc` or `~/.bash_profile`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Then reload:

```bash
source ~/.bashrc
```

#### For Zsh

Add to `~/.zshrc`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Then reload:

```bash
source ~/.zshrc
```

#### For Fish

Add to `~/.config/fish/config.fish`:

```fish
set -gx PATH $PATH (go env GOPATH)/bin
```

Then reload:

```fish
source ~/.config/fish/config.fish
```

## Usage

Navigate to any Go project and run:

```bash
goat
```

The tool will:
- Recursively find all tests in the current directory and subdirectories
- Run `go test -json ./...`
- Display results in an interactive TUI

## Install from GitHub

Once this project is pushed to GitHub, anyone can install it with:

```bash
go install github.com/jesperbjensen/goat@latest
```

Or a specific version:

```bash
go install github.com/jesperbjensen/goat@v1.0.0
```

## Uninstall

To remove goat:

```bash
rm $(go env GOPATH)/bin/goat
```

## Troubleshooting

### "goat: command not found"

- Ensure `$(go env GOPATH)/bin` is in your PATH
- Try running with full path: `$(go env GOPATH)/bin/goat`

### "No tests found"

- Make sure you're in a directory with Go test files
- Verify tests exist with: `go test ./...`
- Test files must end with `_test.go`

### TUI doesn't display correctly

- Ensure your terminal supports ANSI colors
- Try a different terminal (iTerm2, Alacritty, etc.)
- Check terminal size: `echo $COLUMNS x $LINES`

## Development

To run without installing:

```bash
go run .
```

To build locally:

```bash
go build -o goat
./goat
```
