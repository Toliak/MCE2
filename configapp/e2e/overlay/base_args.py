from pathlib import Path
import subprocess
from typing import *


def test_help(binary_e2e: Path, binary_prod: Path) -> bool:
    try:
        output = subprocess.check_output(
            [binary_prod.as_posix(), "-help"], 
            text=True, 
            stderr=subprocess.STDOUT,
        )
    except subprocess.CalledProcessError as e:
        print(f"Binary execution failed: {e}")
        print(e.output)
        return False
    
    print(output)
    cond = (
        'verbosity' in output and
        'no-ui' in output and
        'preset' in output
    )
    return cond

def test_no_ui_just_exit(binary_e2e: Path, binary_prod: Path) -> bool:
    try:
        output = subprocess.check_output(
            [binary_prod.as_posix(), "-no-ui", "-preset={}", "-repo-update-enable=false"], 
            text=True, 
            stderr=subprocess.STDOUT,
        )
    except subprocess.CalledProcessError as e:
        print(f"Binary execution failed: {e}")
        print(e.output)
        return False
    
    print(output)
    
    cond = (
        'verbosity' in output and
        'no-ui' in output and
        'preset' in output
    )
    return cond


def tests_args() -> List[Callable[[Path, Path], bool]]:
    return [
        test_help,
        test_no_ui_just_exit,
    ]