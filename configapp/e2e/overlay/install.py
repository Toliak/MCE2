import json
from pathlib import Path
import subprocess
from typing import *
import util

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
        output = util.check_output_with_live_echo(
            [
                binary_prod.as_posix(), 
                "-preset", 
                json.dumps(preset), 
                "-mce-repo-url=file:///repo",
                "-no-ui",
                "-y",
             ], 
        )
    except subprocess.CalledProcessError as e:
        print(f"Binary execution failed: {e}")
        print(e.output)
        return False
    
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
        output = util.check_output_with_live_echo(
            [
                binary_prod.as_posix(), 
                "-preset", 
                json.dumps(preset), 
                "-mce-repo-url=file:///repo",
                "-no-ui",
                "-y",
             ],
        )
    except subprocess.CalledProcessError as e:
        print(f"Binary execution failed: {e}")
        print(e.output)
        return False
    
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

def test_install_downloads(binary_e2e: Path, binary_prod: Path) -> bool:
    preset = {
        "apps-download": {"en": True},
        "download-lf": {"en": True},
        "download-fzf": {"en": True},
    }
    print("dump:", json.dumps(preset))

    try:
        output = util.check_output_with_live_echo(
            [
                binary_prod.as_posix(), 
                "-preset", 
                json.dumps(preset), 
                "-mce-repo-url=file:///repo",
                "-no-ui",
                "-y",
             ],
        )
    except subprocess.CalledProcessError as e:
        print(f"Binary execution failed: {e}")
        print(e.output)
        return False
    
    lf = (Path.home() / ".local/bin/lf")
    if not lf.exists():
        print("Unable to find lf")
        return False
    
    fzf = (Path.home() / ".local/bin/fzf")
    if not lf.exists():
        print("Unable to find fzf")
        return False
    
    return True

def test_install_zsh_vim_tmux_local(binary_e2e: Path, binary_prod: Path) -> bool:
    preset = {
        "base-cfg-tmux": {"en": True},
        "base-cfg-zsh": {"en": True},
        "bash-config": {"en": True},
        "cfg-local-bash": {"en": True},
        "cfg-local-shared": {"en": True},
        "cfg-local-zsh": {"en": True},
        "cfg-vim-zmix": {"en": True},
        "cfg-zsh-p10k": {"en": True},
        "cfg-zsh-syntax-highlighting": {"en": True},
        "mce2": {"en": True},
        "mce2-repo": {"en": True},
        "os-packages": {"en": True},
        "package-curl": {"en": True},
        "package-git": {"en": True},
        "package-psmisc": {"en": True},
        "package-tmux": {"en": True},
        "package-vim": {"en": True},
        "package-zsh": {"en": True},
        "shared-shell-config": {"en": True},
        "tmux-config": {"en": True},
        "vim-config": {"en": True},
        "zsh-config": {"en": True}
    }
    print("dump:", json.dumps(preset))

    try:
        output = util.check_output_with_live_echo(
            [
                binary_prod.as_posix(), 
                "-preset", 
                json.dumps(preset), 
                "-mce-repo-url=file:///repo",
                "-no-ui",
                "-y",
             ],
        )
    except subprocess.CalledProcessError as e:
        print(f"Binary execution failed: {e}")
        print(e.output)
        return False
    
    zshrc = (Path.home() / ".zshrc")
    print("Zshrc exists: ", zshrc.exists())
    if not zshrc.exists():
        return False

    bashrc = (Path.home() / ".bashrc")
    print("Zshrc exists: ", bashrc.exists())
    if not bashrc.exists():
        return False
    
    zshrc_text = zshrc.read_text()
    if "ZSH_THEME='powerlevel10k/powerlevel10k'" not in zshrc_text:
        print("Unable to find the correct zshrc line")
        return False

    p10k_path = Path.home() / ".local/share/MakeConfigurationEasier2/data/oh-my-zsh/custom/themes/powerlevel10k"
    if not p10k_path.exists():
        print("P10k path does not exist: ", p10k_path)
        return False

    p10k_path = Path.home() / ".local/share/MakeConfigurationEasier2/data/oh-my-zsh/custom/plugins/zsh-syntax-highlighting"
    if not p10k_path.exists():
        print("zsh-syntax-highlighting path does not exist: ", p10k_path)
        return False

    p10k_path = Path.home() / ".local/share/MakeConfigurationEasier2/data/oh-my-tmux"
    if not p10k_path.exists():
        print("oh-my-tmux path does not exist: ", p10k_path)
        return False

    p10k_path = Path.home() / ".local/share/MakeConfigurationEasier2/data/vimrc-amix"
    if not p10k_path.exists():
        print("vimrc-amix path does not exist: ", p10k_path)
        return False
    
    t1 = "\nsource '/home/user/.local/share/MakeConfigurationEasier2/data/local-cfg.bash'"
    bashrc_text = bashrc.read_text()
    if t1 not in bashrc_text:
        print("Unable to find the correct bashrc line: ", t2)
        return False

    t2 = "\nsource '/home/user/.local/share/MakeConfigurationEasier2/data/local-cfg.zsh'"
    zshrc_text = zshrc.read_text()
    if t2 not in zshrc_text:
        print("Unable to find the correct zshrc line: ", t2)
        return False

    t2 = "\nsource '/home/user/.local/share/MakeConfigurationEasier2/data/local-pre-cfg.zsh'"
    zshrc_text = zshrc.read_text()
    if t2 not in zshrc_text:
        print("Unable to find the correct zshrc line: ", t2)
        return False
    
    
    return True

def test_install_p10k_without_oh_my_zsh_must_fail(binary_e2e: Path, binary_prod: Path) -> bool:
    preset = {
        "base-cfg-tmux": {"en": True},
        "cfg-zsh-p10k": {"en": True},
        "mce2": {"en": True},
        "mce2-repo": {"en": True},
        "os-packages": {"en": True},
        "package-curl": {"en": True},
        "package-git": {"en": True},
        "package-psmisc": {"en": True},
        "package-tmux": {"en": True},
        "package-vim": {"en": True},
        "package-zsh": {"en": True},
    }
    print("dump:", json.dumps(preset))

    try:
        output = util.check_output_with_live_echo(
            [
                binary_prod.as_posix(), 
                "-preset", 
                json.dumps(preset), 
                "-mce-repo-url=file:///repo",
                "-no-ui",
                "-y",
             ],
        )
    except subprocess.CalledProcessError as e:
        print(f"Binary execution failed: {e}")
        print(e.output)
        return False
    
    t = "Selected unavailable Tegn 'cfg-zsh-p10k'"
    if t not in output:
        print("Not found text:", t)
        return False

    return True

def tests_install() -> List[Callable[[Path, Path], bool]]:
    return [
        test_install_p10k_without_oh_my_zsh_must_fail,
        test_install_ohmyzsh,
        test_install_zsh_p10k,
        test_install_downloads,
        test_install_zsh_vim_tmux_local,
    ]


# TODO: test the local zsh config