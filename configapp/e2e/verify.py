#!/usr/bin/env python3
import json
import subprocess
import os
import glob
import platform
import shutil
import sys

def check_arch(expected):
    machine = platform.machine()
    # Map platform.machine() to our CPUArch values
    mapping = {
        'x86_64': 'amd64',
        'aarch64': 'aarch64',
        'armv7l': 'armv7',
        'i686': 'i386',
        'i386': 'i386',
        'mips64': 'mips64le',      # approximate
        'ppc64le': 'ppc64',         # approximate
        'riscv64': 'riscv64'
    }
    actual = mapping.get(machine, machine)
    if actual != expected:
        print(f"ARCH mismatch: expected {expected}, got {actual} (machine={machine})")
        return False
    return True

def check_os_type(expected):
    if expected != 'linux':
        print(f"OS type mismatch: expected linux, got {expected}")
        return False
    return True

def check_pkg_manager(raw_cmd):
    path = shutil.which(raw_cmd)
    if path:
        print(f"Package manager '{raw_cmd}' found at '{path}'")
        return True
    else:
        print(f"Package manager '{raw_cmd}' not found in PATH")
        return False

def check_syslib(value, raw):
    if value == "glibc":
        # Look for typical glibc files
        glibc_patterns = [
            "/lib/ld-linux*.so*",
            "/lib/libc.so*",
            "/lib64/ld-linux*.so*",
            "/usr/lib/libc.so*",
            "/lib/*-linux-gnu/libc.so*"
        ]
        for pattern in glibc_patterns:
            if glob.glob(pattern):
                print(f"Found glibc evidence: {pattern}")
                return True
        # Also try ldd --version
        try:
            out = subprocess.check_output(["ldd", "--version"], stderr=subprocess.STDOUT, text=True)
            if "glibc" in out or "GNU libc" in out:
                print("ldd --version indicates glibc")
                return True
        except:
            pass
        print("No glibc evidence found")
        return False

    elif value == "musl":
        musl_patterns = [
            "/lib/ld-musl*.so*",
            "/lib/libc.musl*.so*"
        ]
        for pattern in musl_patterns:
            if glob.glob(pattern):
                print(f"Found musl evidence: {pattern}")
                return True
        if os.path.exists("/etc/alpine-release"):
            print("Found /etc/alpine-release, indicating musl")
            return True
        try:
            out = subprocess.check_output(["ldd", "--version"], stderr=subprocess.STDOUT, text=True)
            if "musl" in out:
                print("ldd --version indicates musl")
                return True
        except:
            pass
        print("No musl evidence found")
        return False

    else:  # "unknown"
        print(f"Unexpected syslib 'unknown' on this platform")
        return False

def check_distrib(expected):
    # expected is a dict with id, id_like, name, version
    os_release = {}
    if os.path.exists("/etc/os-release"):
        with open("/etc/os-release") as f:
            for line in f:
                line = line.strip()
                if line and not line.startswith("#"):
                    if '=' in line:
                        key, val = line.split('=', 1)
                        val = val.strip('"')
                        os_release[key] = val
    else:
        print("/etc/os-release not found")
        return False

    ok = True

    # Compare ID
    if expected.get('id'):
        expected_id = expected['id'].lower()
        actual_id = os_release.get('ID', '').lower()
        if actual_id != expected_id:
            print(f"ID mismatch: expected {expected_id}, got {actual_id}")
            ok = False

    # Compare ID_LIKE (if non‑empty in expected)
    if expected.get('id_like'):
        expected_like = expected['id_like']
        actual_like = os_release.get('ID_LIKE', '').lower().split(' ')

        for id_like in expected_like:
            id_like = id_like.lower()
            
            if id_like not in actual_like:
                print(f"ID_LIKE mismatch: expected {expected_like}, got {actual_like}")
                ok = False

    # Compare NAME
    if expected.get('name'):
        expected_name = expected['name'].lower()
        actual_name = os_release.get('NAME', '').lower()
        if actual_name != expected_name:
            print(f"NAME mismatch: expected {expected_name}, got {actual_name}")
            ok = False

    # Compare VERSION_ID using the raw version string from the JSON
    if expected.get('version', {}).get('raw'):
        expected_ver = expected['version']['raw']
        actual_ver = os_release.get('VERSION_ID', '')
        if actual_ver != expected_ver:
            print(f"VERSION_ID mismatch: expected {expected_ver}, got {actual_ver}")
            ok = False

    return ok

def main():
    binary = "/test/mce"
    if not os.path.exists(binary):
        print(f"Binary {binary} not found")
        sys.exit(1)

    try:
        output = subprocess.check_output([binary, "-harvest-only"], text=True)
    except subprocess.CalledProcessError as e:
        print(f"Binary execution failed: {e}")
        print(e.output)
        sys.exit(1)

    print("Binary output:")
    print(output)

    try:
        data = json.loads(output)
    except json.JSONDecodeError as e:
        print(f"Invalid JSON: {e}")
        sys.exit(1)

    print("Parsed JSON:")
    print(json.dumps(data, indent=2))

    failures = 0

    # Architecture
    arch_val = data.get('arch', {}).get('value')
    if arch_val is None:
        print("Missing 'arch.value'")
        failures += 1
    elif not check_arch(arch_val):
        failures += 1

    # OS type
    ostype_val = data.get('osType', {}).get('value')
    if ostype_val is None:
        print("Missing 'osType.value'")
        failures += 1
    elif not check_os_type(ostype_val):
        failures += 1

    # Package manager
    pkg_raw = data.get('pkgManager', {}).get('raw')
    if pkg_raw is None:
        print("Missing 'pkgManager.raw'")
        failures += 1
    elif not check_pkg_manager(pkg_raw):
        failures += 1

    # System library
    syslib_val = data.get('sysLib', {}).get('value')
    syslib_raw = data.get('sysLib', {}).get('raw')
    if syslib_val is None:
        print("Missing 'sysLib.value'")
        failures += 1
    elif not check_syslib(syslib_val, syslib_raw):
        failures += 1

    # Distribution
    distrib = data.get('distrib')
    if distrib is None:
        print("Missing 'distrib'")
        failures += 1
    elif not check_distrib(distrib):
        failures += 1

    if failures == 0:
        print("\nAll checks passed!")
        sys.exit(0)
    else:
        print(f"\n{failures} check(s) failed.")
        sys.exit(1)

if __name__ == "__main__":
    main()