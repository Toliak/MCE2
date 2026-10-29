package downloads

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
	"github.com/toliak/mce/tegnbuilder"
	"github.com/toliak/mce/tegns/tegn"
)

func calculateGitstatusdInternalName(osInfo *tegnbuilder.OSInfoExt) (string, error) {
	var osPart, archPart string
	switch osInfo.OsType.V {
	case data.OSTypeLinux:
		osPart = "linux"
		switch osInfo.Arch.V {
		case data.CPUArchAMD64:
			archPart = "x86_64"
		case data.CPUArchI386:
			archPart = "i686"
		case data.CPUArchAARCH64:
			archPart = "aarch64"
		case data.CPUArchPPC64:
			archPart = "ppc64le"
		default:
			return "", fmt.Errorf("unsupported CPU arch for gitstatusd: %s, %s", osInfo.OsType.V, osInfo.Arch.V)
		}
	case data.OSTypeDarwin:
		osPart = "darwin"
		switch osInfo.Arch.V {
		case data.CPUArchAMD64:
			archPart = "x86_64"
		case data.CPUArchAARCH64:
			archPart = "arm64"
		default:
			return "", fmt.Errorf("unsupported CPU arch for gitstatusd: %s, %s", osInfo.OsType.V, osInfo.Arch.V)
		}
	// case data.OSTypeWindows: osPart = "msys_nt-10.0"
	default:
		return "", fmt.Errorf("unsupported OS type for gitstatusd: %s", osInfo.OsType.V)
	}
	return fmt.Sprintf("gitstatusd-%s-%s", osPart, archPart), nil
}

// gitstatusdPostProcess is a custom post-processing function for gitstatusd.
// It unpacks the tar.gz, moves the binary to ~/.cache/gitstatus/, and creates a symlink in ~/.local/bin.
func gitstatusdPostProcess(data tegn.GenericDownloadPostProcessData) error {
	// Step 1: Determine the internal name of the binary within the tar.gz
	internalName, err := calculateGitstatusdInternalName(data.OsInfo)
	if err != nil {
		return fmt.Errorf("failed to calculate internal gitstatusd name: %w", err)
	}

	// Step 2: Use the generic tar.gz unpacker to extract the binary to the appDir
	unpacker := tegn.GenericDownloadPostTarGzUnpack(internalName)
	err = unpacker(data) // Pass the original data to the unpacker
	if err != nil {
		return fmt.Errorf("failed to unpack gitstatusd binary: %w", err)
	}

	// At this point, the binary named `internalName` should be in `data.appDir`

	// Step 3: Define paths for the final destination and the symlink
	cacheDir := filepath.Join(data.OsInfo.HomeDir, ".cache", "gitstatus")
	symlinkTargetPath := filepath.Join(cacheDir, internalName) // e.g., ~/.cache/gitstatus/gitstatusd-linux-x86_64
	// WARN: finalBinaryPath must have a name from `calculateGitstatusdInternalName`. Because p10k zsh theme looks for it.
	finalBinaryPath := filepath.Join(data.AppDir, data.AppName) // e.g., ~/.local/bin/gitstatusd

	// Step 4: Ensure the cache directory exists
	err = os.MkdirAll(cacheDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create gitstatus cache directory '%s': %w", cacheDir, err)
	}

	// Step 5: Move the extracted binary from appDir to the cache location
	err = platform.MkdirAllParent(finalBinaryPath)
	if err != nil {
		return fmt.Errorf("failed to create parent directory for binary '%s': %w", finalBinaryPath, err)
	}

	// Step 6: Make the moved binary executable in the cache location
	err = os.Chmod(finalBinaryPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to set executable permissions on '%s': %w", finalBinaryPath, err)
	}

	// Step 7: Create a symlink in the standard appDir location pointing to the cached binary
	// This allows tools expecting 'gitstatusd' in PATH to find it.
	// Ensure the parent directory for the symlink exists (usually ~/.local/bin)
	err = platform.MkdirAllParent(symlinkTargetPath)
	if err != nil {
		return fmt.Errorf("failed to create parent directory for symlink '%s': %w", symlinkTargetPath, err)
	}

	// Create the symlink (relative path preferred if possible, but absolute is safer here)
	err = os.Symlink(finalBinaryPath, symlinkTargetPath)
	if err != nil {
		// If the symlink already exists, try removing it first
		if os.IsExist(err) {
			err = os.Remove(symlinkTargetPath)
			if err != nil {
				return fmt.Errorf("failed to remove existing symlink '%s': %w", symlinkTargetPath, err)
			}
			err = os.Symlink(finalBinaryPath, symlinkTargetPath)
		}
		if err != nil {
			return fmt.Errorf("failed to create symlink '%s' -> '%s': %w", symlinkTargetPath, finalBinaryPath, err)
		}
	}

	return nil
}

// downloadGitstatusd defines the Tegn for downloading and setting up gitstatusd.
var Gitstatusd = tegn.NewGenericDownloadBuilder(
	"gitstatusd",        // appName: Used for the symlink name in ~/.local/bin
	"gitstatusd v1.5.4", // name
	"Fast Git status daemon (gitstatusd binary)", // description
	"gitstatusd.tar.gz",                          // downloadFileName
	func(osInfo *tegnbuilder.OSInfoExt) (string, error) {
		// Define base URL and version based on install.info
		baseURL := "https://github.com/romkatv/gitstatus/releases/download/v1.5.4/"

		internalName, err := calculateGitstatusdInternalName(osInfo)
		if err != nil {
			return "", err
		}

		// Construct the final URL for the .tar.gz release asset
		filename := fmt.Sprintf("%s.tar.gz", internalName)
		url := fmt.Sprintf("%s%s", baseURL, filename)
		return url, nil
	},
	gitstatusdPostProcess, // Use the custom post-processing function
)

// TODO: add a parameter -- add cache symlink or not
