// Copyright (c) 2020, 2022-2024 D. Bohdan
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	tsize "github.com/kopoli/go-terminal-size"
	"github.com/mitchellh/go-wordwrap"
	"github.com/shirou/gopsutil/v4/process"
)

const (
	defaultDumpPath     = ""
	defaultLength       = 20
	defaultMemFormat    = "%.1f"
	defaultNewlines     = false
	defaultOutputPath   = "-"
	defaultQuiet        = false
	defaultRecordTime   = 1000 // ms
	defaultSampleTime   = 200  // ms
	defaultTimeFormat   = "%d:%02d:%04.1f"
	defaultVerbose      = false
	defaultWait         = -1
	sparklineLowMaximum = 10000
	usageDivisor        = 1 << 20 // Report memory usage in binary megabytes.
	version             = "0.8.1"
)

var sparklineTicks = []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

type config struct {
	arguments  []string
	command    string
	dumpPath   string
	length     int
	memFormat  string
	newlines   bool
	outputPath string
	quiet      bool
	record     int
	sample     int
	timeFormat string
	wait       int
}

type MemoryTracker struct {
	timestamps []int64
	values     []int64
	maximum    int64
	mu         sync.RWMutex
}

func (mt *MemoryTracker) AddRecord(timestamp int64, value int64) {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	mt.timestamps = append(mt.timestamps, timestamp)
	mt.values = append(mt.values, value)

	if value > mt.maximum {
		mt.maximum = value
	}
}

func (mt *MemoryTracker) History(count int) ([]int64, []int64, int64) {
	mt.mu.RLock()
	defer mt.mu.RUnlock()

	if count < 0 || count > len(mt.values) {
		count = len(mt.values)
	}

	timestampsCopy := make([]int64, count)
	copy(timestampsCopy, mt.timestamps[len(mt.timestamps)-count:])
	valuesCopy := make([]int64, count)
	copy(valuesCopy, mt.values[len(mt.values)-count:])

	return timestampsCopy, valuesCopy, mt.maximum
}

func wrapForTerm(s string) string {
	size, err := tsize.GetSize()
	if err != nil {
		return s
	}

	return wordwrap.WrapString(s, uint(size.Width))
}

func usage(w io.Writer) {
	s := fmt.Sprintf(
		`Usage: %s [-h] [-v] [-d path] [-l n] [-m fmt] [-n] [-o path] [-q] [-t fmt] [-w ms] [--] command [arg ...]`,
		filepath.Base(os.Args[0]),
	)

	fmt.Fprintln(w, wrapForTerm(s))
}

func help() {
	usage(os.Stdout)

	s := fmt.Sprintf(`
Track the RAM usage (resident set size) of a process and its descendants in real time.

Arguments:
  command
          Command to run

  [arg ...]
          Arguments to the command

Options:
  -h, --help
          Print this help message and exit

  -v, --version
          Print the version number and exit

  -d, --dump path
          File to append full memory usage history to when finished

  -l, --length n
          Sparkline length (default: %d)

  -m, --mem-format fmt
          Format string for memory amounts (default: '%v')

  -n, --newlines
          Print new sparkline on new line instead of over previous

  -o, --output path
          Output file to append to ('%v' for standard error)

  -q, --quiet
          Do not print sparklines, only final report

  -r, --record ms
          How frequently to record/report memory usage in ms (default: %d)

  -s, --sample ms
          How frequently to sample memory usage in ms (default: %d)

  -t, --time-format fmt
          Format string for run time (default: '%v')

  -w, --wait ms
          Set '--sample' and '--record' time simultaneously (that both options override)
`,
		defaultLength,
		defaultMemFormat,
		defaultOutputPath,
		defaultRecordTime,
		defaultSampleTime,
		defaultTimeFormat,
	)

	fmt.Print(wrapForTerm(s))
}

func parseArgs() config {
	cfg := config{
		dumpPath:   defaultDumpPath,
		length:     defaultLength,
		memFormat:  defaultMemFormat,
		outputPath: defaultOutputPath,
		record:     defaultRecordTime,
		sample:     defaultSampleTime,
		timeFormat: defaultTimeFormat,
		wait:       defaultWait,
	}

	usageError := func(message string, badValue interface{}) {
		usage(os.Stderr)
		fmt.Fprintf(os.Stderr, "\nError: "+message+"\n", badValue)
		os.Exit(2)
	}

	// Parse the command-line flags.
	printHelp := false
	printVersion := false

	recondTimeSet := false
	sampleTimeSet := false
	waitTimeSet := false

	var i int
	nextArg := func(flag string) string {
		i++

		if i >= len(os.Args) {
			usageError("no value for option: %s", flag)
		}

		return os.Args[i]
	}

	for i = 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		if arg == "--" {
			i++
			break
		}
		if !strings.HasPrefix(arg, "-") {
			break
		}

		switch arg {

		case "-d", "--dump":
			cfg.dumpPath = nextArg(arg)

		case "-h", "--help":
			printHelp = true

		case "-l", "--length":
			value := nextArg(arg)

			length, err := strconv.Atoi(value)
			if err != nil {
				usageError("invalid length: %v", value)
			}

			cfg.length = length

		case "-m", "--mem-format":
			cfg.memFormat = nextArg(arg)

		case "-n", "--newlines":
			cfg.newlines = true

		case "-o", "--output":
			cfg.outputPath = nextArg(arg)

		case "-q", "--quiet":
			cfg.quiet = true

		case "-r", "--record":
			value := nextArg(arg)
			record, err := strconv.Atoi(value)
			if err != nil {
				usageError("invalid record time: %v", value)
			}

			cfg.record = record
			recondTimeSet = true

		case "-s", "--sample":
			value := nextArg(arg)
			sample, err := strconv.Atoi(value)
			if err != nil {
				usageError("invalid sample time: %v", value)
			}

			cfg.sample = sample
			sampleTimeSet = true

		case "-t", "--time-format":
			cfg.timeFormat = nextArg(arg)

		case "-v", "--version":
			printVersion = true

		case "-w", "--wait":
			value := nextArg(arg)

			wait, err := strconv.Atoi(value)
			if err != nil {
				usageError("invalid wait time: %v", value)
			}

			cfg.wait = wait
			waitTimeSet = true

		default:
			usageError("unknown option: %v", arg)
		}
	}

	if printHelp {
		help()
		os.Exit(0)
	}

	if printVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	// Ensure we have a command.
	if i >= len(os.Args) {
		usageError("command is required%v", "")
	}

	// Set the command and arguments.
	cfg.command = os.Args[i]
	if i+1 < len(os.Args) {
		cfg.arguments = os.Args[i+1:]
	} else {
		cfg.arguments = []string{}
	}

	// Handle the wait option.
	if waitTimeSet {
		if !recondTimeSet {
			cfg.record = cfg.wait
		}

		if !sampleTimeSet {
			cfg.sample = cfg.wait
		}
	}

	return cfg
}

func main() {
	cfg := parseArgs()

	if err := run(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cfg config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	startTime := time.Now().UTC()

	// Prepare stderr or file output.
	output, err := getOutput(cfg.outputPath)
	if err != nil {
		return err
	}
	if output != os.Stderr {
		defer output.Close()
	}

	// We use '\r' to print the sparklines on the same line by default.
	coreFormat := "%s " + cfg.memFormat
	sparklineFormat := "\r" + coreFormat
	if cfg.newlines {
		sparklineFormat = coreFormat + "\n"
	}

	// Start the command.
	cmd := exec.CommandContext(ctx, cfg.command, cfg.arguments...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	// Ensure we shut down the process.
	defer func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}()

	// Get the process.
	proc, err := process.NewProcess(int32(cmd.Process.Pid))
	if err != nil {
		return fmt.Errorf("failed to get process: %w", err)
	}

	// Create the memory-tracking data structures and closures that append to them.
	memTracker := &MemoryTracker{}
	sample := []int64{}

	addSample := func() error {
		mem, err := getMemoryUsage(proc)
		if err != nil {
			return err
		}

		sample = append(sample, int64(mem))

		return nil
	}

	addRecord := func() {
		if len(sample) == 0 {
			return
		}

		memTracker.AddRecord(time.Now().UnixNano(), slices.Max(sample))

		if !cfg.quiet {
			_, values, maximum := memTracker.History(cfg.length)
			line := sparkline(maximum, values)
			fmt.Fprintf(output, sparklineFormat, line, float64(maximum)/usageDivisor)
		}

		sample = []int64{}
	}

	// Start memory tracking by adding an initial record before we wait.
	_ = addSample()
	addRecord()

	done := make(chan error, 1)

	go func() {
		sampleTicker := time.NewTicker(time.Duration(cfg.sample) * time.Millisecond)
		defer sampleTicker.Stop()

		recordTicker := time.NewTicker(time.Duration(cfg.record) * time.Millisecond)
		defer recordTicker.Stop()

		for {
			select {

			case <-ctx.Done():
				return

			case <-sampleTicker.C:
				err := addSample()
				if err != nil {
					continue
				}

			case <-recordTicker.C:
				addRecord()
			}
		}
	}()

	go func() {
		done <- cmd.Wait()
	}()

	// Wait for either the command's completion or a signal.
	select {

	case err := <-done:
		// Stop memory tracking.
		cancel()

		if err != nil {
			return err
		}

	case sig := <-sigChan:
		cancel()

		return fmt.Errorf("received signal: %v", sig)
	}

	// Get the complete final stats.
	timestamps, values, maximum := memTracker.History(-1)
	endTime := time.Now().UTC()

	if len(values) == 0 {
		fmt.Fprintln(output, "no data collected")
	} else {
		if !cfg.newlines && !cfg.quiet {
			fmt.Fprintln(output)
		}

		summary := summarize(values, maximum, startTime, endTime, cfg.memFormat, cfg.timeFormat)
		fmt.Fprintln(output, summary)
	}

	// Dump the memory usage history if required.
	if cfg.dumpPath != defaultDumpPath {
		if err := dumpHistory(cfg.dumpPath, timestamps, values); err != nil {
			return fmt.Errorf("failed to dump history: %w", err)
		}
	}

	return nil
}

func getOutput(path string) (*os.File, error) {
	if path == defaultOutputPath {
		return os.Stderr, nil
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open output file: %w", err)
	}

	return file, nil
}

func getMemoryUsage(proc *process.Process) (int64, error) {
	children, err := proc.Children()
	if err != nil && err != process.ErrorNoChildren {
		return 0, err
	}

	var total uint64
	for _, child := range children {
		mem, err := child.MemoryInfo()

		// If we can't get memory info for a child, skip it.
		if err == nil && mem != nil && mem.RSS > 0 {
			total += mem.RSS
		}
	}

	memInfo, err := proc.MemoryInfo()
	if err != nil {
		return 0, fmt.Errorf("failed to get process memory info: %w", err)
	}
	if memInfo == nil {
		return 0, fmt.Errorf("no memory info available")
	}

	total += memInfo.RSS
	return int64(total), nil
}

func summarize(values []int64, maximum int64, start, end time.Time, memFormat, timeFormat string) string {
	avg := average(values)

	result := strings.Builder{}

	result.WriteString(" avg: ")
	result.WriteString(fmt.Sprintf(memFormat, float64(avg)/usageDivisor))
	result.WriteString("\n max: ")
	result.WriteString(fmt.Sprintf(memFormat, float64(maximum)/usageDivisor))
	result.WriteString("\ntime: ")
	hours, minutes, seconds := hmsDelta(start, end)
	result.WriteString(fmt.Sprintf(timeFormat, hours, minutes, seconds))

	return result.String()
}

func average[T int64](values []T) T {
	var sum T
	for _, value := range values {
		sum += value
	}

	if len(values) == 0 {
		return T(0)
	}

	return sum / T(len(values))
}

func dumpHistory(path string, timestamps []int64, values []int64) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for i, timestamp := range timestamps {
		_, err := fmt.Fprintf(writer, "%d %d\n", timestamp/1_000_000, values[i])
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func hmsDelta(start, end time.Time) (int, int, float64) {
	delta := end.Sub(start)
	totalMillis := int(delta / time.Millisecond)

	hours := totalMillis / (60 * 60 * 1000)
	remaining := totalMillis % (60 * 60 * 1000)
	minutes := remaining / (60 * 1000)
	remaining = remaining % (60 * 1000)
	seconds := float64(remaining) / 1000.0

	return hours, minutes, seconds
}

func sparkline(maximum int64, data []int64) string {
	if maximum <= sparklineLowMaximum {
		return strings.Repeat(string(sparklineTicks[0]), max(1, len(data)))
	}

	tickMax := int64(len(sparklineTicks) - 1)
	result := strings.Builder{}

	for _, x := range data {
		tickIndex := int(tickMax * x / maximum)
		result.WriteRune(sparklineTicks[tickIndex])
	}

	return result.String()
}
