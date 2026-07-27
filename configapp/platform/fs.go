package platform

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// Checks if a command is available in PATH
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !(err != nil && os.IsNotExist(err))
}

// Mkdir with parents
func MkdirAllParent(path string) error {
	parentDir := filepath.Dir(path)
	return os.MkdirAll(parentDir, 0755)
}

// See https://opensource.com/article/18/6/copying-files-go
func CopyFile(src, dst string) error {
	// TODO: preserve permissions see copyAndDeleteFile
	sourceFileStat, err := os.Stat(src)
	if err != nil {
			return err
	}

	if !sourceFileStat.Mode().IsRegular() {
			return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
			return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
			return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func copyAndDeleteFile(src, dst string, perm os.FileMode) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer sourceFile.Close()

	// Create/overwrite destination preserving permissions
	// similar to os.Create
	destFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("create dest: %w", err)
	}

	// Stream copy bytes safely
	// TODO: use `CopyFile`
	_, err = io.Copy(destFile, sourceFile)
	closeErr := destFile.Close()

	if err != nil {
		os.Remove(dst) // Clean up partial destination file on copy failure
		return fmt.Errorf("copy content: %w", err)
	}
	if closeErr != nil {
		os.Remove(dst)
		return fmt.Errorf("close dest file: %w", closeErr)
	}

	// Explicitly close source before removing it (required for Windows)
	sourceFile.Close()

	// 5. Remove original file after successful copy
	if err := os.Remove(src); err != nil {
		return fmt.Errorf("remove original file: %w", err)
	}

	return nil
}

// Helper to detect EXDEV errors across Linux, macOS, and Windows
func isCrossDeviceError(err error) bool {
	var linkErr *os.LinkError
	if errors.As(err, &linkErr) {
		if errno, ok := linkErr.Err.(syscall.Errno); ok {
			// EXDEV: Linux/macOS/Windows cross-device error
			return errno == syscall.EXDEV
		}
	}
	return false
}

// MoveFile moves a single file from src to dst across any storage devices or filesystems.
func MoveFile(src, dst string) error {
	// 1. Ensure the source path points to a regular file (not a directory)
	info, err := os.Lstat(src)
	if err != nil {
		return fmt.Errorf("stat source file: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("source '%s' is a directory, not a file", src)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("source '%s' is not a regular file (e.g. socket, pipe, device)", src)
	}

	// 2. Attempt fast, atomic move on the same filesystem
	err = os.Rename(src, dst)
	if err == nil {
		return nil // Success on same device!
	}

	// 3. Check if error is due to cross-device move
	if !isCrossDeviceError(err) {
		return fmt.Errorf("rename failed: %w", err)
	}

	// 4. Fallback: Copy and delete across devices
	return copyAndDeleteFile(src, dst, info.Mode().Perm())
}

func OpenOrCreateFileAppend(path string) (*os.File, error) {
	if FileExists(path) {
		v, err := os.OpenFile(path, os.O_APPEND | os.O_WRONLY, 0644)
		return v, err
	}

	v, err := os.Create(path)
	return v, err
}

// Safe function to append text into the file.
// It will create the file if it does not exist.
func AppendFilepathString(path string, text string) error {
	outputFile, err := OpenOrCreateFileAppend(path)
    if err != nil {
        return err
    }
    defer outputFile.Close()

	_, err = outputFile.WriteString(text)
	if err != nil {
        return err
    }

	return nil
}

func IsFileSymlinkTo(pathToCheck string, symlinkPointsAt string) bool {
	// Readlink returns the raw string destination of the symbolic link
	target, err := os.Readlink(pathToCheck)
	if err != nil {
		// Returns an error if the file is not a symlink or doesn't exist
		return false 
	}

	targetAbs, err := filepath.Abs(target)
	if err != nil {
		return false 
	}

	symlinkPointsAtAbs, err := filepath.Abs(symlinkPointsAt)
	if err != nil {
		return false 
	}
	
	// Check for an exact match to your target string
	return targetAbs == symlinkPointsAtAbs
}
