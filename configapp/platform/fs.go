package platform

import (
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"
)

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

func copyAndDeleteFile(src, dst string, perm os.FileMode) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer sourceFile.Close()

	// Create/overwrite destination preserving permissions
	destFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("create dest: %w", err)
	}

	// Stream copy bytes safely
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