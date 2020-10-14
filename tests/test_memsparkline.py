#! /usr/bin/env python3

# Copyright (c) 2020 D. Bohdan
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
import os.path
import subprocess
import unittest


TEST_PATH = os.path.dirname(os.path.realpath(__file__))
COMMAND = os.environ.get(
    'MEMSPARKLINE_COMMAND',
    os.path.join(TEST_PATH, '..', 'memsparkline')
)


def run(*args, check=True, stdin=None):
    return subprocess.run(
        [COMMAND, *args],
        check=check,
        stdin=stdin,
        stderr=subprocess.PIPE,
        stdout=subprocess.PIPE,
    )


class TestMemsparkline(unittest.TestCase):

    def test_usage(self):
        self.assertRegex(
            run(check=False).stderr.decode('ascii'),
            r'^usage',
        )

    def test_basic(self):
        self.assertRegex(
            run('sleep', '1').stderr.decode('utf-8'),
            r'(?s).*avg:.*max:',
        )

    def test_length(self):
        stderr = run('-l', '10', '-w', '10', 'sleep', '1') \
            .stderr.decode('utf-8')

        self.assertRegex(
            stderr,
            r'(?m)\r[^ ]{10} \d+\.\d{2}\navg',
        )

    def test_wait_1(self):
        stderr = run('-w', '2000', 'sleep', '1').stderr.decode('utf-8')

        self.assertEqual(len(stderr.split('\n')), 4)

    def test_wait_2(self):
        stderr = run('-n', '-w', '100', 'sleep', '1').stderr.decode('utf-8')

        self.assertIn(len(stderr.split('\n')), range(10, 15))


if __name__ == '__main__':
    unittest.main()
