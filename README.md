# memsparkline

Track the RAM usage ([resident set size](https://en.wikipedia.org/wiki/Resident_set_size)) of a process, its children, its children's children, etc. in real time with a Unicode text [sparkline](https://en.wikipedia.org/wiki/Sparkline).
See the average and the maximum usage after the process exits, as well as the run time.

## Examples

```none
> memsparkline -- chromium-browser --incognito http://localhost:8081/
▁▁▁▁▄▇▇▇█ 789.5
 avg: 371.0
 max: 789.5
time: 0:00:12.0
```

```none
> memsparkline -n -o log du /usr/ >/dev/null 2>&1 &
> tail -f log
█ 2.8
▆█ 3.3
▆▇█ 3.6
▆▇▇█ 3.9
▆▇▇█▆ 3.3
▆▇▇█▆▆ 3.3
▆▇▇█▆▆▆ 3.3
▆▇▇█▆▆▆▆ 3.3
▄▅▅▆▅▅▅▅█ 5.2
▄▅▅▆▅▅▅▅██ 5.2
 avg: 3.7
 max: 5.2
time: 0:00:10.1
```

## Installation

### Prebuilt binaries

Prebuilt binaries for
FreeBSD (amd64),
Linux (aarch64, riscv64, x86_64),
macOS (arm64, x86_64),
OpenBSD (amd64),
and Windows (amd64, x86)
are attached to [releases](https://github.com/dbohdan/memsparkline/releases).

### Go

Install Go, then run the following command:

```shell
go install github.com/dbohdan/memsparkline@latest
```

## Build requirements

- Go 1.21
- OS supported by [gopsutil](https://github.com/shirou/gopsutil)
- POSIX Make for testing

## Compatibility and limitations

memsparkline works on POSIX systems supported by [gopsutil](https://github.com/shirou/gopsutil).
It has been tested on Debian, Ubuntu, FreeBSD, and OpenBSD.
Unfortunately, gopsutil doesn't support NetBSD.
NetBSD users can install the last [Python release](https://pypi.org/project/memsparkline/) of memsparkline.

Although memsparkline seems to work on Windows, Windows support has received little testing outside of [CI](https://en.wikipedia.org/wiki/Continuous_integration).
The sparkline displays incorrectly in the Command Prompt and [ConEmu](https://conemu.github.io/) on Windows 7 with the stock console fonts.
It displays correctly on Windows 10 with the font NSimSun.

## Operation

### Usage

```none
Usage: memsparkline [-h] [-v] [-d path] [-l n] [-m fmt] [-n] [-o path] [-q] [-t
fmt] [-w ms] [--] command [arg ...]

Track the RAM usage (resident set size) of a process and its descendants in
real time.

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
          Sparkline length (default: 20)

  -m, --mem-format fmt
          Format string for memory amounts (default: '%.1f')

  -n, --newlines
          Print new sparkline on new line instead of over previous

  -o, --output path
          Output file to append to ('-' for standard error)

  -q, --quiet
          Do not print sparklines, only final report

  -r, --record ms
          How frequently to record/report memory usage in ms (default: 1000)

  -s, --sample ms
          How frequently to sample memory usage in ms (default: 200)

  -t, --time-format fmt
          Format string for run time (default: '%d:%02d:%04.1f')

  -w, --wait ms
          Set '--sample' and '--record' time simultaneously (that both options
override)
```

### Samples and records

memsparkline differentiates between _samples_ and _records_.
Samples are measurements of memory usage.
Records are information about memory usage printed to the chosen output (given by `--output`) and added to history (saved using the `--dump` option).

There is a separate setting for the sample time and the record time.
The sample time determines the interval between when memory usage is measured.
The record time determines the interval between when a record is made (written to the output and added to history).
When sampling is more frequent than recording (as with the default settings),
memsparkline uses the highest sampled value since the last record.

A short sample time like 5 ms can result in high CPU usage,
up to 100% of one CPU core.
To reduce CPU usage, sample less frequently.
The default sample time of 200 ms results in memsparkline using around 15% of a 2019 x86-64 core on the developer's machine.

Records are only created when at least one sample has been taken.
Setting the record time shorter than the sample time is allowed for convenience, but no record is added when there are no samples.

## License

MIT.

## See also

memusg and spark (both linked below) inspired this project.

### Tracking memory usage

* [DragonFly BSD](https://man.dragonflybsd.org/?command=time&section=ANY),
  [FreeBSD](https://man.freebsd.org/cgi/man.cgi?query=time&format=html),
  [macOS](https://ss64.com/osx/time.html),
  [NetBSD](https://man.netbsd.org/time.1),
  and [OpenBSD](https://man.openbsd.org/time)
  time(1) flag `-l`.
* [GNU time(1)](https://linux.die.net/man/1/time) flag `-v`.
* [memusg](http://gist.github.com/526585) — a Bash script for FreeBSD, Linux, and macOS that measures the peak resident set size of a process.

### Sparklines

* [spark](https://github.com/holman/spark) — a Bash script that generates a Unicode text sparkline from a list of numbers.
* [sparkline.tcl](https://wiki.tcl-lang.org/page/Sparkline) — a Tcl script inspired by spark made by the developer of this project.
  Adds a `--min` and `--max` option for setting the scale.
