package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

var (
	command  = getCommand()
	testPath = getCurrentDir()
)

func getCommand() []string {
	if envCmd := os.Getenv("MEMSPARKLINE_COMMAND"); envCmd != "" {
		return strings.Fields(envCmd)
	}

	return []string{"./memsparkline"}
}

func getCurrentDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
}

func runMemsparkline(t *testing.T, args ...string) (string, string, error) {
	// Start with command args, if any, then add the test-specific args.
	allArgs := append([]string{}, command[1:]...)
	allArgs = append(allArgs, args...)

	cmd := exec.Command(command[0], allArgs...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func getSleepCommand(duration float64) []string {
	return []string{"test/sleep", strconv.FormatFloat(duration, 'f', -1, 64)}
}

func TestUsage(t *testing.T) {
	_, stderr, _ := runMemsparkline(t)
	if matched, _ := regexp.MatchString("^Usage", stderr); !matched {
		t.Error("Expected usage information in stderr")
	}
}

func TestVersion(t *testing.T) {
	stdout, _, _ := runMemsparkline(t, "-v")
	if matched, _ := regexp.MatchString(`\d+\.\d+\.\d+`, stdout); !matched {
		t.Error("Expected version number in stdout")
	}
}

func TestUnknownOptBeforeHelp(t *testing.T) {
	_, _, err := runMemsparkline(t, "--foo", "--help")

	if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() != 2 {
		t.Errorf("Expected exit status 2, got %v", err)
	}
}

func TestUnknownOptAfterHelp(t *testing.T) {
	_, _, err := runMemsparkline(t, "--help", "--foo")

	if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() != 2 {
		t.Errorf("Expected exit status 2, got %v", err)
	}
}

func TestBasic(t *testing.T) {
	args := getSleepCommand(0.5)
	_, stderr, _ := runMemsparkline(t, args...)
	if matched, _ := regexp.MatchString(`(?s).*avg:.*max:`, stderr); !matched {
		t.Error("Expected 'avg:' and 'max:' in output")
	}
}

func TestEndOfOptions(t *testing.T) {
	args := append([]string{"--"}, getSleepCommand(0.1)...)

	_, stderr, _ := runMemsparkline(t, args...)
	if matched, _ := regexp.MatchString(`(?s).*avg:.*max:`, stderr); !matched {
		t.Error("Expected 'avg:' and 'max:' in output")
	}
}

func TestEndOfOptionsHelp(t *testing.T) {
	args := append([]string{"--"}, getSleepCommand(0.1)...)
	args = append(args, "-h")

	_, stderr, _ := runMemsparkline(t, args...)
	if matched, _ := regexp.MatchString(`(?s).*avg:.*max:`, stderr); !matched {
		t.Error("Expected 'avg:' and 'max:' in output")
	}
}

func TestLength(t *testing.T) {
	args := append([]string{"-l", "5", "-w", "10"}, getSleepCommand(0.5)...)
	_, stderr, _ := runMemsparkline(t, args...)

	if matched, _ := regexp.MatchString(`(?m)\r[^ ]{5} \d+\.\d\r?\n avg`, stderr); !matched {
		t.Error("Expected sparkline of specific length followed by summary")
	}
}

func TestMemFormat(t *testing.T) {
	args := append([]string{"-l", "5", "-w", "10", "-m", "%0.2f"}, getSleepCommand(0.5)...)
	_, stderr, _ := runMemsparkline(t, args...)

	if matched, _ := regexp.MatchString(`(?m)\r[^ ]{5} \d+\.\d{2}\r?\n avg`, stderr); !matched {
		t.Error("Expected sparkline with memory format with two decimal places")
	}
}

func TestTimeFormat(t *testing.T) {
	args := append([]string{"-l", "10", "-t", "%d:%05d:%06.3f"}, getSleepCommand(0.5)...)
	_, stderr, _ := runMemsparkline(t, args...)

	if matched, _ := regexp.MatchString(`(?m)time: \d+:\d{5}:\d{2}\.\d{3}\r?\n`, stderr); !matched {
		t.Error("Expected specific time format in summary")
	}
}

func TestWait1(t *testing.T) {
	args := append([]string{"-w", "2000"}, getSleepCommand(0.5)...)
	_, stderr, _ := runMemsparkline(t, args...)

	if lines := strings.Count(stderr, "\n"); lines != 4 {
		t.Errorf("Expected 4 lines in output, got %d", lines)
	}
}

func TestWait2(t *testing.T) {
	args := append([]string{"-n", "-w", "10"}, getSleepCommand(0.5)...)
	_, stderr, _ := runMemsparkline(t, args...)

	if lines := strings.Count(stderr, "\n"); lines < 9 {
		t.Errorf("Expected at least 9 lines in output, got %d", lines)
	}
}

func TestSampleAndRecord(t *testing.T) {
	args := append([]string{"-r", "500", "-s", "100"}, getSleepCommand(0.5)...)
	_, stderr, _ := runMemsparkline(t, args...)

	if lines := strings.Count(stderr, "\n"); lines != 4 {
		t.Errorf("Expected 4 lines in output, got %d", lines)
	}
}

func TestQuiet(t *testing.T) {
	args := append([]string{"-q"}, getSleepCommand(0.5)...)
	_, stderr, _ := runMemsparkline(t, args...)

	if matched, _ := regexp.MatchString("^ avg", stderr); !matched {
		t.Error("Expected output to start with 'avg' in quiet mode")
	}
}

func TestMissingBinary(t *testing.T) {
	_, stderr, err := runMemsparkline(t, "no-such-binary-exists")

	if err == nil {
		t.Error("Expected error for missing binary")
	}

	if !strings.Contains(stderr, "failed to start command") {
		t.Error("Expected 'failed to start command' in stderr")
	}
}

func TestDoubleDash(t *testing.T) {
	stdout, _, _ := runMemsparkline(t, "--", "ls", "-l")

	if !strings.Contains(stdout, "\n") {
		t.Error("Expected newline in output")
	}
}

func TestTwoDoubleDashes(t *testing.T) {
	args := append(command, "--", "ls", "-l")
	stdout, _, _ := runMemsparkline(t, append([]string{"--"}, args...)...)

	if !strings.Contains(stdout, "\n") {
		t.Error("Expected newline in output")
	}
}

func TestDump(t *testing.T) {
	dumpPath := filepath.Join(testPath, "dump.log")
	// Clean up any existing file so we don't append to it.
	os.Remove(dumpPath)

	args := append([]string{"-q", "-w", "100", "-d", dumpPath}, getSleepCommand(0.5)...)
	runMemsparkline(t, args...)

	content, err := os.ReadFile(dumpPath)
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) == 0 {
		t.Error("Expected non-empty dump file")
	}

	for _, line := range lines {
		if matched, _ := regexp.MatchString(`\d+ \d+`, line); !matched {
			t.Errorf("Invalid line format: %s", line)
		}
	}
}

func TestOutput(t *testing.T) {
	outputPath := filepath.Join(testPath, "output.log")
	// Clean up any existing file so we don't append to it.
	os.Remove(outputPath)

	args := append([]string{"-q", "-o", outputPath}, getSleepCommand(0.5)...)
	for i := 0; i < 2; i++ {
		runMemsparkline(t, args...)
	}

	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) != 7 {
		t.Errorf("Expected 7 lines in output file, got %d", len(lines))
	}
}
