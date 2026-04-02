import os
from pathlib import Path
import sys
from typing import *
from os_info import test_os_info
from base_args import tests_args
from install import tests_install

test_list: List[Callable[[Path, Path], bool]] = [
    *tests_install(),
    test_os_info,
    *tests_args(),
]

def main():
    binary = Path("/test/bin/mce-e2e")
    binary_prod = Path("/test/bin/mce")
    if not binary.exists():
        print(f"Binary {binary} not found")
        return False
    if not binary_prod.exists():
        print(f"Binary prod {binary} not found")
        return False

    tests_len = len(test_list)

    failed_tests: List[Callable[[Path], bool]] = []
    for i, v in enumerate(test_list):
        print(f"[{i+1}/{tests_len}] TEST {v.__name__}")
        r = v(binary, binary_prod)
        if r:
            print(f"[{i+1}/{tests_len}] TEST {v.__name__} OK")
        else:
            failed_tests.append(v)
            print(f"[{i+1}/{tests_len}] TEST {v.__name__} FAIL")

    if failed_tests:
        print("Failed tests:")
        for v in failed_tests:
            print(f"- {v.__name__}")

        sys.exit(1)

    print(f"All tests ({tests_len}) passed")


if __name__ == "__main__":
    main()