# Changelog

All notable changes to this project will be documented in this file.

## [0.0.2] - 2024-01-18

### Added
- Command-line argument support - pass any `go test` flags or paths
- Comprehensive INSTALL.md with installation instructions
- .gitignore for build artifacts
- Support for custom test paths (e.g., `goat ./pkg/mypackage`)
- Support for go test flags (e.g., `goat -v -race -cover ./...`)

### Changed
- Default behavior now uses `./...` to recursively find all tests
- Refactored codebase into modular structure:
  - `types.go` - Core data structures
  - `styles.go` - UI styling definitions
  - `ui.go` - Bubble Tea TUI implementation
  - `test_runner.go` - Test execution and parsing
  - `editor.go` - Editor integration
  - `main.go` - Application entry point
- Updated README with command-line usage examples

### Fixed
- Tests now discoverable in subdirectories when run from project root
- Binary no longer tracked in git

## [0.0.1] - 2024-01-17

### Added
- Initial release
- Interactive TUI for viewing Go test results
- Sidebar with test list (failed tests shown first)
- Content pane with detailed test output
- Keyboard navigation (arrow keys, j/k, g/G)
- Mouse support (click to select, scroll wheel)
- Press Enter to open test file at error line in editor
- Color-coded test results (green pass, red fail)
- Support for multiple editors (Zed, VS Code, Vim, etc.)
