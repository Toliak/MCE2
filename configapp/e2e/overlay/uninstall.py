import json
from pathlib import Path
import re
import subprocess
from typing import *
import util

def test_uninstall_not_installed(binary_e2e: Path, binary_prod: Path) -> bool:
    preset = ["not-installed-1"]
    try:
        output = util.check_output_with_live_echo(
            [
                binary_prod.as_posix(), 
                "uninstall", 
                "-preset", 
                json.dumps(preset), 
                "-repo-update-enable=0",
                "-no-ui",
                "-y",
             ], 
        )
    except subprocess.CalledProcessError as e:
        print(f"Binary execution failed: {e}")
        print(e.output)
        return False
    
    if re.match('^Removed not-installed-1.+ it is not installed', output, re.M):
        return False
    
    return True

def tests_uninstall() -> List[Callable[[Path, Path], bool]]:
    return [
        test_uninstall_not_installed,
    ]


# TODO: test the local zsh config


# TODO: test the blocker tegns somehow????
