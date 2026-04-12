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
    if not zshrc.exists():
        return False
    
    zshrc_text = zshrc.read_text()
    if "export ZSH='/home/user/.local/share/MakeConfigurationEasier2/data/oh-my-zsh'" not in zshrc_text:
        print("Unable to find the correct ZSH line")
        return False
    
    return True

def test_install_zsh_p10k(binary_e2e: Path, binary_prod: Path) -> bool:
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
        "cfg-zsh-p10k": {"en": True},
    }
    print("dump:", json.dumps(preset))

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
    if not zshrc.exists():
        return False
    
    zshrc_text = zshrc.read_text()
    if "ZSH_THEME='powerlevel10k/powerlevel10k'" not in zshrc_text:
        print("Unable to find the correct ZSH line")
        return False

    p10k_path = Path.home() / ".local/share/MakeConfigurationEasier2/data/oh-my-zsh/custom/themes/powerlevel10k"
    if not p10k_path.exists():
        print("P10k path does not exist: ", p10k_path)
        return False
    
    return True

def tests_install() -> List[Callable[[Path, Path], bool]]:
    return [
        test_install_ohmyzsh,
        test_install_zsh_p10k
    ]


# TODO: test the local zsh config