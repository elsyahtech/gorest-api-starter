package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type testEvent struct {
	Action string
	Output string
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	root, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("resolve repository directory: %w", err)
	}
	reportDirectory := filepath.Join(root, "coverage")
	if err := os.MkdirAll(reportDirectory, 0o755); err != nil {
		return fmt.Errorf("create integration report directory: %w", err)
	}
	coverageDirectory, err := os.MkdirTemp("", "gorest-integration-coverage-")
	if err != nil {
		return fmt.Errorf("create temporary integration coverage directory: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(coverageDirectory); err != nil {
			fmt.Fprintf(os.Stderr, "remove temporary coverage data: %v\n", err)
		}
	}()

	jsonPath := filepath.Join(reportDirectory, "integration-test.json")
	jsonReport, err := os.Create(jsonPath)
	if err != nil {
		return fmt.Errorf("create JSON integration report: %w", err)
	}
	command := exec.Command(
		"go", "test", "-tags=integration", "-v", "-json", "-count=1", "-timeout=5m", "./integration-test",
	)
	command.Env = setEnvironment(os.Environ(), "GOREST_INTEGRATION_COVERAGE_DIR", coverageDirectory)
	command.Stdout = io.MultiWriter(os.Stdout, jsonReport)
	command.Stderr = os.Stderr
	testErr := command.Run()
	closeErr := jsonReport.Close()
	if closeErr != nil {
		return fmt.Errorf("close JSON integration report: %w", closeErr)
	}

	textPath := filepath.Join(reportDirectory, "integration-test.txt")
	if err := writeReadableReport(jsonPath, textPath); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "Integration reports: %s and %s\n", jsonPath, textPath)

	entries, err := os.ReadDir(coverageDirectory)
	if err != nil {
		return fmt.Errorf("read integration coverage data: %w", err)
	}
	if len(entries) == 0 {
		if testErr != nil {
			return fmt.Errorf("integration tests failed: %w", testErr)
		}
		return errors.New("integration tests passed but produced no application coverage data")
	}

	profilePath := filepath.Join(reportDirectory, "integration.coverprofile")
	if err := runTool("go", "tool", "covdata", "textfmt", "-i="+coverageDirectory, "-o="+profilePath); err != nil {
		return fmt.Errorf("convert application integration coverage: %w", err)
	}
	functionPath := filepath.Join(reportDirectory, "integration-coverage.txt")
	functionReport, err := os.Create(functionPath)
	if err != nil {
		return fmt.Errorf("create integration coverage summary: %w", err)
	}
	cover := exec.Command("go", "tool", "cover", "-func="+profilePath)
	cover.Stdout = functionReport
	cover.Stderr = os.Stderr
	coverErr := cover.Run()
	closeErr = functionReport.Close()
	if coverErr != nil {
		return fmt.Errorf("summarize application integration coverage: %w", coverErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close integration coverage summary: %w", closeErr)
	}
	htmlPath := filepath.Join(reportDirectory, "integration-coverage.html")
	if err := runTool("go", "tool", "cover", "-html="+profilePath, "-o="+htmlPath); err != nil {
		return fmt.Errorf("generate application integration coverage HTML: %w", err)
	}
	fmt.Fprintf(os.Stdout, "Application coverage: %s, %s, and %s\n", profilePath, functionPath, htmlPath)
	if testErr != nil {
		return fmt.Errorf("integration tests failed: %w", testErr)
	}
	return nil
}

func writeReadableReport(jsonPath, textPath string) error {
	input, err := os.Open(jsonPath)
	if err != nil {
		return fmt.Errorf("open JSON integration report: %w", err)
	}
	defer input.Close()
	output, err := os.Create(textPath)
	if err != nil {
		return fmt.Errorf("create readable integration report: %w", err)
	}

	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)
	for scanner.Scan() {
		var event testEvent
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			_ = output.Close()
			return fmt.Errorf("decode Go test JSON event: %w", err)
		}
		if event.Output != "" {
			if _, err := io.WriteString(output, event.Output); err != nil {
				_ = output.Close()
				return fmt.Errorf("write readable integration report: %w", err)
			}
		}
	}
	scanErr := scanner.Err()
	closeErr := output.Close()
	if scanErr != nil {
		return fmt.Errorf("read JSON integration report: %w", scanErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close readable integration report: %w", closeErr)
	}
	return nil
}

func runTool(name string, args ...string) error {
	command := exec.Command(name, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func setEnvironment(environment []string, name, value string) []string {
	prefix := name + "="
	result := make([]string, 0, len(environment)+1)
	for _, entry := range environment {
		if len(entry) >= len(prefix) && entry[:len(prefix)] == prefix {
			continue
		}
		result = append(result, entry)
	}
	return append(result, prefix+value)
}
