package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// loadTests runs go test -json and parses the results
func loadTests() tea.Msg {
	cmd := exec.Command("go", "test", "-json")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return testsLoadedMsg{tests: []TestResult{}}
	}

	if err := cmd.Start(); err != nil {
		return testsLoadedMsg{tests: []TestResult{}}
	}

	scanner := bufio.NewScanner(stdout)
	testResults := make(map[string]*TestResult)

	for scanner.Scan() {
		line := scanner.Bytes()
		var event TestEvent

		if err := json.Unmarshal(line, &event); err != nil {
			continue
		}

		// Handle test output events
		if event.Test != "" && event.Action == "output" {
			if testResults[event.Test] == nil {
				testResults[event.Test] = &TestResult{
					Name:   event.Test,
					Output: []string{},
				}
			}
			testResults[event.Test].Output = append(testResults[event.Test].Output, event.Output)

			// Parse file path and line number from output (e.g., "    main_test.go:6: Fake error")
			if testResults[event.Test].FilePath == "" {
				trimmed := strings.TrimSpace(event.Output)
				if strings.Contains(trimmed, ".go:") {
					parts := strings.Split(trimmed, ":")
					if len(parts) >= 2 {
						testResults[event.Test].FilePath = parts[0]
						if len(parts) >= 2 {
							var lineNum int
							fmt.Sscanf(parts[1], "%d", &lineNum)
							testResults[event.Test].Line = lineNum
						}
					}
				}
			}
		}

		// Handle test pass/fail events
		if event.Test != "" && (event.Action == "pass" || event.Action == "fail") {
			if testResults[event.Test] == nil {
				testResults[event.Test] = &TestResult{
					Name:   event.Test,
					Output: []string{},
				}
			}
			testResults[event.Test].Status = event.Action
		}
	}

	cmd.Wait()

	// Convert map to sorted slice
	tests := make([]TestResult, 0, len(testResults))
	for _, result := range testResults {
		tests = append(tests, *result)
	}

	// Sort: failed tests first, then by name
	sort.Slice(tests, func(i, j int) bool {
		if tests[i].Status == "fail" && tests[j].Status != "fail" {
			return true
		}
		if tests[i].Status != "fail" && tests[j].Status == "fail" {
			return false
		}
		return tests[i].Name < tests[j].Name
	})

	return testsLoadedMsg{tests: tests}
}
