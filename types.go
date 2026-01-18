package main

// TestEvent represents a single event from go test -json output
type TestEvent struct {
	Time    string  `json:"Time"`
	Action  string  `json:"Action"`
	Package string  `json:"Package"`
	Test    string  `json:"Test"`
	Elapsed float64 `json:"Elapsed"`
	Output  string  `json:"Output"`
}

// TestResult represents the aggregated result of a test
type TestResult struct {
	Name     string
	Status   string
	Output   []string
	FilePath string
	Line     int
}

// testsLoadedMsg is sent when tests have been loaded and parsed
type testsLoadedMsg struct {
	tests []TestResult
}
