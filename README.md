# GOAT - Go Automated Testing TUI

A Terminal User Interface (TUI) for running and viewing Go tests with an interactive sidebar and detailed output view.

## Features

- 🎯 Interactive TUI with keyboard and mouse navigation
- ✅ Color-coded test results (green for pass, red for fail)
- 📝 Detailed test output view
- 🔗 Clickable file links to open tests in your editor
- 🎨 Beautiful terminal UI using Bubble Tea and Lip Gloss
- 🚀 Failed tests displayed first for quick debugging

## Installation

```bash
go build -o goat
```

Or install as a Go tool:

```bash
go install
```

## Usage

Run the test viewer in any Go project directory:

```bash
goat
```

Or with a local build:

```bash
./goat
```

### Command-Line Arguments

You can pass any `go test` arguments to goat:

```bash
# Run tests in a specific package
goat ./pkg/mypackage

# Run tests with verbose output
goat -v ./...

# Run specific tests by pattern
goat -run TestMyFunction ./...

# Run tests with coverage
goat -cover ./...

# Run tests with race detector
goat -race ./...

# Combine multiple flags
goat -v -race -cover ./pkg/...
```

**Default behavior**: If no arguments are provided, goat runs `go test -json ./...` (all tests recursively).

## Key Bindings

- **↑/k** - Move up in test list
- **↓/j** - Move down in test list
- **g/Home** - Jump to first test
- **G/End** - Jump to last test
- **Enter** - Open test file at error line in editor
- **q/Ctrl+C** - Quit application

## Mouse Controls

- **Click** - Select a test from the sidebar
- **Scroll Wheel** - Navigate up/down through tests

## Project Structure

The codebase is organized into focused modules for easy maintenance:

### `main.go`

Entry point of the application. Initializes and runs the Bubble Tea program.

### `types.go`

Core data structures:

- `TestEvent` - Individual events from `go test -json` output
- `TestResult` - Aggregated test result with status and output
- `testsLoadedMsg` - Message type for when tests are loaded

### `styles.go`

UI styling definitions using Lip Gloss:

- `sidebarStyle` - Test list sidebar styling
- `contentStyle` - Main content area styling
- `selectedItemStyle` - Currently selected test styling
- `passStyle` - Green styling for passing tests
- `failStyle` - Red styling for failing tests
- `titleStyle` - Section title styling

### `test_runner.go`

Test execution and parsing:

- `loadTests()` - Runs `go test -json` and parses output
- Extracts file paths and line numbers from test errors
- Sorts tests with failed tests first

### `editor.go`

Editor integration:

- `openFile()` - Opens files at specific lines in available editors
- Supports: Zed, VS Code, Sublime Text, Atom, Vim, Neovim, Emacs
- Automatically detects which editor is available

### `ui.go`

Bubble Tea TUI implementation:

- `model` - Application state
- `Init()` - Initialization
- `Update()` - Event handling (keyboard, mouse, messages)
- `View()` - Rendering the UI
- `renderSidebar()` - Builds the test list sidebar
- `renderContent()` - Builds the detailed test output view
- `calculateVisibleRange()` - Handles scrolling for long test lists

## How It Works

1. Application parses command-line arguments (defaults to `./...` if none provided)
2. Runs `go test -json [your args...]`
3. Tests are discovered based on your arguments
4. Test output is parsed line-by-line as JSON events
5. Test results are aggregated and sorted (failures first)
6. TUI displays tests in a sidebar with status indicators
7. Selecting a test shows detailed output in the content pane
8. Pressing Enter opens the test file at the error line

## Supported Editors

The application will attempt to open files in these editors (in order):

1. Zed (`zed`)
2. VS Code (`code`)
3. Sublime Text (`subl`)
4. Atom (`atom`)
5. Vim (`vim`)
6. Neovim (`nvim`)
7. Emacs (`emacs`)

## Requirements

- Go 1.25.6 or later
- A terminal that supports ANSI colors
- (Optional) One of the supported editors for the "open file" feature

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Terminal styling

## License

MIT
