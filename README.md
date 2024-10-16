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


## Compatibility and limitations

memsparkline works on POSIX systems supported by [psutil](https://github.com/giampaolo/psutil).
It has been tested on Debian, Ubuntu, FreeBSD, NetBSD, and OpenBSD.

Although memsparkline seems to work on Windows, Windows support has received little testing outside of [CI](https://en.wikipedia.org/wiki/Continuous_integration).
The sparkline displays incorrectly in the Command Prompt and [ConEmu](https://conemu.github.io/) on Windows 7 with the stock console fonts.
It displays correctly on Windows 10 with the font NSimSun.


## Operation

### Usage

```none
usage: memsparkline [-h] [-v] [-d path] [-l n] [-m fmt] [-n] [-o path] [-q]
                    [-r ms] [-s ms] [-t fmt] [-w ms]
                    command ...

Track the RAM usage (resident set size) of a process and its descendants in
real time.

positional arguments:
  command               command to run
  args                  arguments to command

options:
  -h, --help            show this help message and exit
  -v, --version         show program's version number and exit
  -d path, --dump path  file in which to write full memory usage history when
                        finished
  -l n, --length n      sparkline length (default: 20)
  -m fmt, --mem-format fmt
                        format string for memory amounts (default: "%0.1f")
  -n, --newlines        print new sparkline on new line instead of over
                        previous
  -o path, --output path
                        output file to append to ("-" for standard error)
  -q, --quiet           do not print sparklines, only final report
  -r ms, --record ms    how frequently to record/report memory usage (default:
                        every 1000 ms)
  -s ms, --sample ms    how frequently to sample memory usage (default: every
                        200 ms)
  -t fmt, --time-format fmt
                        format string for run time (default: "%d:%02d:%04.1f")
  -w ms, --wait ms      set "--sample" and "--record" time simultaneously
                        (that both options override)
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
The default sample time of 200 ms results in memsparkline using around 10% of a 2019 x86-64 core on the developer's machine.

Records are only created after a sample has been taken.
Setting the record time shorter than the sample time is allowed for convenience but equivalent to setting it to the sample time.


## Installation

memsparkline requires Python 3.8 or later.

### Installing from PyPI

The recommended way to install memsparkline is [from PyPI](https://pypi.org/project/memsparkline/) with [pipx](https://github.com/pypa/pipx).

```sh
pipx install memsparkline
```

You can also use pip:

```sh
pip install --user memsparkline
```

### Manual installation

1. Install the dependencies from the package repositories for your OS.
   You will find instructions for some operating systems below.
2. Download `src/memsparkline/main.py` and copy it to a directory in `PATH` as `memsparkline`.
   For example:

```sh
git clone https://github.com/dbohdan/memsparkline
sudo install memsparkline/src/memsparkline/main.py /usr/local/bin/memsparkline
```

#### Dependencies

##### Debian/Ubuntu

```sh
sudo apt install python3-psutil
```

##### DragonFly BSD 6.6, FreeBSD 13.1

```sh
sudo pkg install py39-psutil
```

##### NetBSD 9.3

```sh
sudo pkgin in py310-psutil
```

##### OpenBSD

```sh
doas pkg_add py3-psutil
```


## License

MIT.


## See also

memusg and spark (both linked below) inspired this project.

### Tracking memory usage

* [DragonFly BSD](https://man.dragonflybsd.org/?command=time&section=ANY), [FreeBSD](https://man.freebsd.org/cgi/man.cgi?query=time&format=html), [NetBSD](https://man.netbsd.org/time.1), [OpenBSD](https://man.openbsd.org/time), and [macOS](https://ss64.com/osx/time.html) time(1) flag `-l`.
* [GNU time(1)](https://linux.die.net/man/1/time) flag `-v`.
* [memusg](http://gist.github.com/526585) — a Bash script for FreeBSD, Linux, and macOS that measures the peak resident set size of a process.

### Sparklines

* [spark](https://github.com/holman/spark) — a Bash script that generates a Unicode text sparkline from a list of numbers.
* [sparkline.tcl](https://wiki.tcl-lang.org/page/Sparkline) — a Tcl script by the developer of this project that does the same. Adds a `--min` and `--max` option for setting the scale.
