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

import os
import re
import shlex
import subprocess
import sys
import unittest
from pathlib import Path

TEST_PATH = Path(__file__).resolve().parent
COMMAND = shlex.split(os.environ.get("MEMSPARKLINE_COMMAND", ""))
if [] == COMMAND:
    COMMAND = [sys.executable, Path(TEST_PATH, "..", "memsparkline.py")]


def run(
    *args: str,
    check: bool = True,
    return_stdout: bool = False,
    return_stderr: bool = True,
) -> str:
    completed = subprocess.run(
        COMMAND + list(args),  # noqa: S603
        check=check,
        stdin=None,
        capture_output=True,
    )

    output = ""
    if return_stdout:
        output += completed.stdout.decode("utf-8")
    if return_stderr:
        output += completed.stderr.decode("utf-8")

    return output


class TestMemsparkline(unittest.TestCase):
    def test_usage(self) -> None:
        assert re.search("^usage", run(check=False))

    def test_version(self) -> None:
        assert re.search("\\d+\\.\\d+\\.\\d+", run("-v", return_stdout=True))


@unittest.skipUnless(os.name == "posix", "requires a POSIX OS")
class TestMemsparklinePOSIX(unittest.TestCase):
    def test_basic(self) -> None:
        assert re.search("(?s).*avg:.*max:", run("sleep", "1"))

    def test_length(self) -> None:
        stderr = run("-l", "10", "-w", "10", "sleep", "1")

        assert re.search("(?m)\\r[^ ]{10} \\d+\\.\\d\\n avg", stderr)

    def test_mem_format(self) -> None:
        stderr = run("-l", "10", "-w", "10", "-m", "%0.2f", "sleep", "1")

        assert re.search("(?m)\\r[^ ]{10} \\d+\\.\\d{2}\\n avg", stderr)

    def test_time_format(self) -> None:
        stderr = run("-l", "10", "-t", "%d:%05d:%06.3f", "sleep", "1")

        assert re.search("(?m)time: \\d+\\:\\d{5}:\\d{2}\\.\\d{3}\\n", stderr)

    def test_wait_1(self) -> None:
        stderr = run("-w", "2000", "sleep", "1")

        assert len(stderr.split("\n")) == 5

    def test_wait_2(self) -> None:
        stderr = run("-n", "-w", "100", "sleep", "1")

        assert len(stderr.split("\n")) in range(10, 15)

    def test_quiet(self) -> None:
        stderr = run("-q", "sleep", "1")

        assert re.search("^ avg", stderr)

    def test_missing_binary(self) -> None:
        with self.assertRaises(subprocess.CalledProcessError) as err:
            run("no-such-binary-exists")
            assert re.search(
                r"No such file or directory",
                err.exception.output,
            )

    def test_double_dash(self) -> None:
        assert "\n" in run("--", "ls", "-l", return_stdout=True)

    def test_two_double_dashes(self) -> None:
        assert "\n" in run("--", *COMMAND, "--", "ls", "-l", return_stdout=True)

    def test_output(self) -> None:
        output_path = Path(TEST_PATH, "output.log")
        if output_path.exists():
            output_path.unlink()

        for _ in range(2):
            run("-q", "-o", str(output_path), "sleep", "1")

        text = output_path.read_text()
        assert len(text.split("\n")) == 7


@unittest.skipUnless(os.name == "nt", "requires Windows")
class TestMemsparklineWindows(unittest.TestCase):
    def test_cmd_basic(self) -> None:
        assert re.search("(?s).*avg:.*max:", run("cmd.exe", "/c", "dir"))

    def test_cmd_pause(self) -> None:
        stderr = run("-w", "2000", "cmd", "/c", "timeout", "/t", "1")

        assert len(stderr.split("\n")) == 5


if __name__ == "__main__":
    unittest.main()
