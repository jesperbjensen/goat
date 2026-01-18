package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

const (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
)

type TestEvent struct {
	Time    string  `json:"Time"`
	Action  string  `json:"Action"`
	Package string  `json:"Package"`
	Test    string  `json:"Test"`
	Elapsed float64 `json:"Elapsed"`
	Output  string  `json:"Output"`
}

type TestResult struct {
	Status   string
	Output   []string
	FilePath string
	Line     int
}

func main() {
	// Check for --only-fail flag
	onlyFail := slices.Contains(os.Args[1:], "--only-fail")

	// Run go test with JSON output
	cmd := exec.Command("go", "test", "-json")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating stdout pipe: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting go test: %v\n", err)
		os.Exit(1)
	}

	// Parse JSON output line by line
	scanner := bufio.NewScanner(stdout)
	testResults := make(map[string]*TestResult)

	for scanner.Scan() {
		line := scanner.Bytes()
		var event TestEvent

		if err := json.Unmarshal(line, &event); err != nil {
			continue
		}

		// Track test output
		if event.Test != "" && event.Action == "output" {
			if testResults[event.Test] == nil {
				testResults[event.Test] = &TestResult{Output: []string{}}
			}
			testResults[event.Test].Output = append(testResults[event.Test].Output, event.Output)

			// Parse file path and line number from output (e.g., "    main_test.go:6: Fake error")
			if testResults[event.Test].FilePath == "" {
				trimmed := strings.TrimSpace(event.Output)
				if strings.Contains(trimmed, ".go:") {
					parts := strings.Split(trimmed, ":")
					if len(parts) >= 2 {
						testResults[event.Test].FilePath = parts[0]
						// Parse line number
						if len(parts) >= 2 {
							var lineNum int
							fmt.Sscanf(parts[1], "%d", &lineNum)
							testResults[event.Test].Line = lineNum
						}
					}
				}
			}
		}

		// We only care about pass/fail actions for individual tests
		if event.Test != "" && (event.Action == "pass" || event.Action == "fail") {
			if testResults[event.Test] == nil {
				testResults[event.Test] = &TestResult{Output: []string{}}
			}
			testResults[event.Test].Status = event.Action
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading test output: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.Wait(); err != nil {
		// Test failures cause non-zero exit, but we still want to print results
		// so we don't exit here
	}

	// Print results in the requested format
	for test, result := range testResults {
		status := result.Status

		// Skip passing tests if --only-fail is set
		if onlyFail && status == "pass" {
			continue
		}

		statusUpper := ""
		color := ""
		if status == "pass" {
			statusUpper = "PASS"
			color = colorGreen
		} else if status == "fail" {
			statusUpper = "FAIL"
			color = colorRed
		}

		fmt.Printf("* %s: %s%s%s", test, color, statusUpper, colorReset)

		// Add clickable link for failed tests
		if status == "fail" && result.FilePath != "" && result.Line > 0 {
			fmt.Printf(" (%s:%d)", result.FilePath, result.Line)
		}
		fmt.Println()

		// Show output for failed tests
		if status == "fail" && len(result.Output) > 0 {
			for _, line := range result.Output {
				// Trim whitespace and skip empty lines
				trimmed := strings.TrimSpace(line)
				if trimmed != "" {
					fmt.Printf("  %s\n", line)
				}
			}
		}
	}
}
