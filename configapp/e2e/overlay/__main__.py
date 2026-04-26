import os
from pathlib import Path
import sys
from typing import *
from uninstall import tests_uninstall
from os_info import test_os_info
from base_args import tests_args
from install import tests_install
import argparse
from dataclasses import dataclass

TestFunT = Callable[[Path, Path], bool]

test_list: List[TestFunT] = [
    *tests_uninstall(),
    *tests_install(),
    test_os_info,
    *tests_args(),
]

@dataclass
class Args:
    selected_tests: Optional[List[str]] = None
    show_list: bool = False


def parse_args(args: List[str]) -> Args:
    parser = argparse.ArgumentParser(description="Run tests")
    parser.add_argument(
        "--selected-tests",
        "-s",
        nargs="+",
        help="Run only selected tests by function name",
        default=None
    )
    parser.add_argument(
        "--show-list",
        "-l",
        action="store_true",
        help="Show all available tests"
    )

    args = parser.parse_args(args)
    return Args(selected_tests=args.selected_tests, show_list=args.show_list)


def main():
    argv = sys.argv[1:]
    args = parse_args(argv)

    if args.show_list:
        print("Available tests:")
        for test_func in test_list:
            print(f"- {test_func.__name__}")
        return

    binary = Path("/test/bin/mce-e2e")
    binary_prod = Path("/test/bin/mce")
    if not binary.exists():
        print(f"Binary {binary} not found")
        return False
    if not binary_prod.exists():
        print(f"Binary prod {binary_prod} not found")
        return False

    # Filter test list based on selected tests
    if args.selected_tests:
        filtered_tests: List[TestFunT] = []
        not_found_tests: List[str] = []
        for sel_test in args.selected_tests:
            for t in test_list:
                if t.__name__ == sel_test:
                    filtered_tests.append(t)
                    break
            else:
                not_found_tests.append(sel_test)

        if not_found_tests:
            print("Not found test names:")
            for test_name in not_found_tests:
                print(f"- {test_name}")
            return False

        test_list_to_run = filtered_tests
    else:
        test_list_to_run = test_list

    tests_len = len(test_list_to_run)

    failed_tests: List[Callable[[Path], bool]] = []
    for i, v in enumerate(test_list_to_run):
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

    print(f"All tests (selected) ({tests_len}) passed")


if __name__ == "__main__":
    main()