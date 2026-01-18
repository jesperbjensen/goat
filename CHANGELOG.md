# Changelog

All notable changes to this project will be documented in this file.

## [0.0.4] - 2024-01-18

### Added

- Scrollable content pane for long test output (Ctrl+D/PgDn to scroll down, Ctrl+U/PgUp to scroll up)
- Scroll indicators show when content extends beyond visible area
- Smart text wrapping for long test names in sidebar
- Intelligent line wrapping for test output in content pane

### Changed

- Content pane now properly handles output longer than screen height
- Test names wrap across multiple lines instead of truncating
- Scroll position resets when changing selected test
- Improved text wrapping algorithm that respects word boundaries

### Fixed

- Overflow issues when test output is very long
- Test names no longer cut off awkwardly in sidebar
- Content pane properly displays all output with scrolling
- Line wrapping now works correctly for long lines

## [0.0.3] - 2024-01-18

### Added

- Loading status displayed while tests are running
- Error handling with clear error messages if test execution fails
- Filter toggle ('f' key) to show only failed tests
- Success message when all tests pass (visible in filter mode)
- Status bar at bottom showing:
  - Pass/fail counts
  - Keyboard shortcuts (f: toggle filter, q: quit)
- New styles for loading, error, success, and dimmed text

### Changed

- Sidebar now shows "Tests (Failures Only)" when filter is active
- Content pane shows success celebration when all tests pass
- UI layout adjusted to accommodate status bar
- Improved user feedback throughout the testing process

### Fixed

- Cursor resets to 0 when toggling filter to prevent out-of-bounds errors
- Proper error propagation from test runner to UI

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
