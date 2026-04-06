import json
from pathlib import Path
import subprocess
from typing import *

def test_install_ohmyzsh(binary_e2e: Path, binary_prod: Path) -> bool:
    preset = {
        "mce2": {"en": True},
        "mce2-repo": {"en": True},
        "os-packages": {"en": True},
        "package-git": {"en": True},
        "package-psmisc": {"en": True},
        "package-tmux": {"en": True},
        "package-zsh": {"en": True},
        "zsh-config": {"en": True},
        "base-cfg-zsh": {"en": True},
    }
    try:
        output = subprocess.check_output(
            [
                binary_prod.as_posix(), 
                "-preset", 
                json.dumps(preset), 
                "-mce-repo-url=file:///repo",
                "-no-ui",
                "-y",
             ], 
            text=True, 
            stderr=subprocess.STDOUT,
        )
    except subprocess.CalledProcessError as e:
        print(f"Binary execution failed: {e}")
        print(e.output)
        return False
    
    print(output)
    zshrc = (Path.home() / ".zshrc")
    print("Zshrc exists: ", zshrc.exists())
    cond = zshrc.exists()
    return cond

def tests_install() -> List[Callable[[Path, Path], bool]]:
    return [
        test_install_ohmyzsh,
    ]