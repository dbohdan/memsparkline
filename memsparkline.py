#! /usr/bin/env python3

# Copyright (c) 2020, 2022-2023 D. Bohdan
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.

import argparse
import contextlib
from datetime import datetime
import sys
import time
from typing import Iterator, IO, List, Tuple

import psutil

SPARKLINE_TICKS = ["▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"]
USAGE_DIVISOR = 1 << 20
VERSION = "0.3.0"


def main() -> None:
    args = cli(sys.argv)

    with open_output(args.output_path, sys.stderr) as output:  # type: IO[str]
        try:
            start_dt = datetime.now()
            process = psutil.Popen([args.command] + args.arguments)
            maximum, history = track(
                process,
                output,
                newlines=args.newlines,
                sparkline_length=args.length,
                wait=args.wait,
                mem_format=args.mem_format,
                quiet=args.quiet,
            )
            process.wait()

            if history == []:
                print("no data collected", file=output)
            else:
                if not args.newlines and not args.quiet:
                    print("", file=output)
                summary = summarize(
                    history,
                    maximum,
                    start_dt,
                    datetime.now(),
                    args.mem_format,
                    args.time_format,
                )
                print("\n".join(summary), file=output)

            if args.dump_path != "":
                with open(args.dump_path, "w") as hist_file:
                    for value in history:
                        print(value, file=hist_file)
        except Exception as err:
            print(
                "error: %s" % err,
                file=output,
            )
            sys.exit(1)

        sys.exit(process.returncode)


def hms_delta(
    start_dt: datetime,
    end_dt: datetime,
) -> Tuple[int, int, float]:
    delta = end_dt - start_dt
    total_millis = (
        delta.days * 24 * 60 * 60 * 1000
        + delta.seconds * 1000
        + delta.microseconds // 1000
    )

    hours, rem = divmod(total_millis, 60 * 60 * 1000)
    minutes, rem = divmod(rem, 60 * 1000)
    seconds = rem / 1000

    return hours, minutes, seconds


def summarize(
    history: List[int],
    maximum: int,
    start_dt: datetime,
    end_dt: datetime,
    mem_format: str,
    time_format: str,
) -> List[str]:
    summary = []
    summary.append(
        (" avg: " + mem_format) % (sum(history) / len(history) / USAGE_DIVISOR)
    )
    summary.append((" max: " + mem_format) % (maximum / USAGE_DIVISOR))

    summary.append("time: " + time_format % hms_delta(start_dt, end_dt))

    return summary


def cli(argv: List[str]) -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Track the RAM usage (resident set size) of a process and "
        "its descendants in real time.",
    )
    parser.add_argument(
        "command",
        default=[],
        help="command to run",
    )
    parser.add_argument(
        "arguments",
        default=[],
        help="arguments to command",
        metavar="arg",
        nargs="*",
    )
    parser.add_argument(
        "-v",
        "--version",
        action="version",
        version=VERSION,
    )
    parser.add_argument(
        "-d",
        "--dump",
        default="",
        dest="dump_path",
        help="file to which to write full memory usage history when finished",
        metavar="path",
    )
    parser.add_argument(
        "-l",
        "--length",
        default=20,
        dest="length",
        help="sparkline length (default: %(default)d)",
        metavar="n",
        type=int,
    )
    parser.add_argument(
        "-m",
        "--mem-format",
        default="%0.1f",
        dest="mem_format",
        help="format string for memory amounts (default: %(default)s)",
        metavar="fmt",
        type=str,
    )
    parser.add_argument(
        "-n",
        "--newlines",
        action="store_true",
        default=False,
        help="print new sparkline on new line instead of over previous",
    )
    parser.add_argument(
        "-o",
        "--output",
        default="-",
        dest="output_path",
        help='output file ("-" for standard error)',
        metavar="path",
    )
    parser.add_argument(
        "-q",
        "--quiet",
        dest="quiet",
        action="store_true",
        help="do not print sparklines, only final report",
    )
    parser.add_argument(
        "-t",
        "--time-format",
        default="%d:%02d:%04.1f",
        dest="time_format",
        help="format string for run time (default: %(default)s)",
        metavar="fmt",
        type=str,
    )
    parser.add_argument(
        "-w",
        "--wait",
        default=1000,
        dest="wait",
        help="how long to wait between taking samples (default: %(default)d)",
        metavar="ms",
        type=int,
    )

    args = parser.parse_args(argv[1:])

    return args


@contextlib.contextmanager
def open_output(path: str, fallback: IO[str]) -> Iterator[IO[str]]:
    handle = fallback
    if path != "-":
        handle = open(path, "w", 1)

    try:
        yield handle
    finally:
        if handle is not sys.stderr:
            handle.close()


def track(
    parent: psutil.Popen,
    output: IO[str],
    newlines: bool = False,
    sparkline_length: int = 20,
    wait: int = 1000,
    mem_format: str = "0.1f%",
    quiet: bool = False,
) -> Tuple[int, List[int]]:
    core_fmt = "%s " + mem_format
    fmt = core_fmt + "\n" if newlines else "\r" + core_fmt
    history = []
    maximum = 0

    try:
        while parent.is_running() and parent.status() != psutil.STATUS_ZOMBIE:
            tree = parent.children(recursive=True)
            tree.append(parent)

            total = sum(x.memory_info().rss for x in tree)
            maximum = max(maximum, total)
            history.append(total)

            if not quiet:
                latest = history[-sparkline_length:]
                line = sparkline(0, maximum, latest)
                print(
                    fmt % (line, total / USAGE_DIVISOR),
                    end="",
                    file=output,
                )

            time.sleep(wait / 1000)
    except KeyboardInterrupt:
        pass

    return (maximum, history)


def sparkline(minimum: float, maximum: float, data: List[float]) -> str:
    tick_max = len(SPARKLINE_TICKS) - 1

    if minimum == maximum:
        return SPARKLINE_TICKS[tick_max // 2] * len(data)

    return "".join(
        SPARKLINE_TICKS[int(tick_max * (x - minimum) / (maximum - minimum))]
        for x in data
    )


if __name__ == "__main__":
    main()
