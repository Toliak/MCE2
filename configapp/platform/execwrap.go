package platform

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

type ExecCommandWrapper struct {
	BufferStdout       bool
	BufferStderr       bool
	RetransmitStdout   bool
	RetransmitStderr   bool
	ThrowExitCodeError bool
	AdditionalEnv      []string
	NeedsRoot          bool
}

type ExecCommandOption func(*ExecCommandWrapper)

func WithBufferStdout(v bool) ExecCommandOption {
	return func(w *ExecCommandWrapper) {
		w.BufferStdout = v
	}
}

func WithBufferStderr(v bool) ExecCommandOption {
	return func(w *ExecCommandWrapper) {
		w.BufferStderr = v
	}
}

func WithRetransmitStdout(v bool) ExecCommandOption {
	return func(w *ExecCommandWrapper) {
		w.RetransmitStdout = v
	}
}

func WithRetransmitStderr(v bool) ExecCommandOption {
	return func(w *ExecCommandWrapper) {
		w.RetransmitStderr = v
	}
}

func WithThrowExitCodeError(v bool) ExecCommandOption {
	return func(w *ExecCommandWrapper) {
		w.ThrowExitCodeError = v
	}
}

func WithAdditionalEnv(env string) ExecCommandOption {
	return func(w *ExecCommandWrapper) {
		w.AdditionalEnv = append(w.AdditionalEnv, env)
	}
}

func WithAdditionalEnvList(env []string) ExecCommandOption {
	return func(w *ExecCommandWrapper) {
		w.AdditionalEnv = append(w.AdditionalEnv, env...)
	}
}

func WithNeedsRoot(v bool) ExecCommandOption {
	return func(w *ExecCommandWrapper) {
		w.NeedsRoot = v
	}
}

func NewExecCommandWrapper(opts ...ExecCommandOption) ExecCommandWrapper {
	w := ExecCommandWrapper{
		BufferStdout:       false,
		BufferStderr:       false,
		RetransmitStdout:   true,
		RetransmitStderr:   true,
		ThrowExitCodeError: false,
		AdditionalEnv:      make([]string, 0),
		NeedsRoot:          false,
	}
	for _, opt := range opts {
		opt(&w)
	}
	return w
}

// threadSafeBuffer wraps bytes.Buffer to allow concurrent writes
type threadSafeBuffer struct {
	mu  sync.Mutex
	buf *bytes.Buffer
}

func (tsb *threadSafeBuffer) Write(p []byte) (n int, err error) {
	tsb.mu.Lock()
	defer tsb.mu.Unlock()
	return tsb.buf.Write(p)
}

func execCommandInternal(name string, arg ...string) *exec.Cmd {
	// logs
	fmt.Printf("Exec: %s %v\n", name, arg)
	return exec.Command(name, arg...)
}

func ExecCommand(config ExecCommandWrapper, name string, arg ...string) (*bytes.Buffer, error) {
	var rawBuffer bytes.Buffer

	// Wrap it for thread safety if both stdout and stderr are buffering
	var bufferWriter io.Writer = &rawBuffer
	
	if config.BufferStdout && config.BufferStderr {
		// If both are writing to the same buffer, we need a mutex
		bufferWriter = &threadSafeBuffer{buf: &rawBuffer}
	}

	usesWrapperForPrivileges := false
	
	var cmd *exec.Cmd
	if !config.NeedsRoot || IsRoot() {
		cmd = execCommandInternal(name, arg...)
	} else {
		usesWrapperForPrivileges = true
		if CommandExists("sudo") {
			newArgs := []string{"-E"}
			newArgs = append(newArgs, name)
			newArgs = append(newArgs, arg...)
			cmd = execCommandInternal("sudo", newArgs...)
		} else if CommandExists("pkexec") {
			newArgs := []string{"--keep-cwd"}
			newArgs = append(newArgs, name)
			newArgs = append(newArgs, arg...)
			cmd = execCommandInternal("pkexec", newArgs...)
		} else {
			return nil, fmt.Errorf("Unable to find app to run as root (sudo or pkexec)")
		}
	}

	cmd.Env = append(os.Environ(), config.AdditionalEnv...)

	// Configure stdout writer(s)
	var stdoutWriters []io.Writer
	if config.RetransmitStdout {
		stdoutWriters = append(stdoutWriters, os.Stdout)
	}
	if config.BufferStdout {
		stdoutWriters = append(stdoutWriters, bufferWriter)
		if usesWrapperForPrivileges {
			// WARNING
			fmt.Println("BufferStdout: Privilege evaluation was used. The sudo/pkexec prompt will be buffered too")
		}
	}
	if len(stdoutWriters) == 0 {
		stdoutWriters = []io.Writer{io.Discard}
	}

	// Configure stderr writer(s)
	var stderrWriters []io.Writer
	if config.RetransmitStderr {
		stderrWriters = append(stderrWriters, os.Stdout)
	}
	if config.BufferStderr {
		stderrWriters = append(stderrWriters, bufferWriter) // reuse buffer if needed
		if usesWrapperForPrivileges {
			// WARNING
			fmt.Println("BufferStderr: Privilege evaluation was used. The sudo/pkexec prompt will be buffered too")
		}
	}
	if len(stderrWriters) == 0 {
		stderrWriters = []io.Writer{io.Discard}
	}

	// Create pipes
	cmd.Stdout = io.MultiWriter(stdoutWriters...)
	cmd.Stderr = io.MultiWriter(stderrWriters...)

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("ExecCommand error start: %w", err)
	}

	cmdErr := cmd.Wait()

	// Wait for the command to finish
	if cmdErr != nil {
		if config.ThrowExitCodeError {
			return nil, fmt.Errorf("ExecCommand error command: %w", cmdErr)
		}
		var exitError *exec.ExitError
		if errors.As(cmdErr, &exitError) {
			return &rawBuffer, nil
		}
		return nil, fmt.Errorf("ExecCommand error command: %w", cmdErr)
	}

	return &rawBuffer, nil
}
