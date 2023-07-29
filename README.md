# memsparkline

Track the RAM usage ([resident set size](https://en.wikipedia.org/wiki/Resident_set_size)) of a process, its children, its children's children, etc. in real time with a Unicode text [sparkline](https://en.wikipedia.org/wiki/Sparkline). See the average and the maximum usage after the process exits, as well as the run time.


## Examples

```none
> memsparkline -- chromium-browser --incognito http://localhost:8081/
▁▁▁▁▄▇▇▇█ 789.5
 avg: 371.0
 max: 789.5
time: 0:00:12.3
```

```none
> memsparkline -o foo command &
> tail -f foo
```


## Installation

memsparkline requires Python 3.7 or later.

### Installing from PyPI

The recommended way to install memsparkline is with [pipx](https://gitlab.com/dbohdan/memsparkline).

```sh
pipx install memsparkline
```

You can also use pip:

```sh
pip install --user memsparkline
```

### Manual installation

1. Install the dependencies using the OS-specific instrutions below.
2. Download `memsparkline.py` and copy it to a directory in `PATH`. For example,

```sh
git clone https://gitlab.com/dbohdan/memsparkline
cd memsparkline
sudo install memsparkline.py /usr/local/bin/memsparkline
```

How to install the dependencies:

#### Debian/Ubuntu

```sh
sudo apt install python3-psuti
```

#### DragonFly BSD 6.6, FreeBSD 13.1

```sh
sudo pkg install py39-psutil
```

#### NetBSD 9.3

```sh
sudo pkgin in py310-psutil
```

#### OpenBSD

```sh
doas pkg_add py3-psutil
```


## Usage

```none
usage: memsparkline.py [-h] [-d path] [-m fmt] [-t fmt] [-l n] [-n] [-o path]
                       [-v] [-w ms]
                       command [arg [arg ...]]

Track the RAM usage (resident set size) of a process and its descendants in
real time.

positional arguments:
  command               command to run
  arg                   arguments to command

optional arguments:
  -h, --help            show this help message and exit
  -d path, --dump path  file in which to write full memory usage history when
                        finished
  -m fmt, --mem-format fmt
                        format string for memory numbers (default: %0.1f)
  -t fmt, --time-format fmt
                        format string for run time (default: %d:%02d:%04.1f)
  -l n, --length n      sparkline length (default: 20)
  -n, --newlines        print new sparkline on new line instead of over
                        previous
  -o path, --output path
                        output file ("" or "-" for standard error)
  -v, --version         show program's version number and exit
  -w ms, --wait ms      how long to wait between taking samples (default:
                        1000)
```


## License

MIT.
