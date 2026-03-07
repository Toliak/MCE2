package platform

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/toliak/mce/osinfo"
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

func execCommandInternal(name string, arg ...string) *exec.Cmd {
	// logs
	fmt.Printf("Exec: %s %v\n", name, arg)
	return exec.Command(name, arg...)
}

func ExecCommand(config ExecCommandWrapper, name string, arg ...string) (*bytes.Buffer, error) {
	var stdoutBuf bytes.Buffer

	usesWrapperForPrivileges := false
	
	var cmd *exec.Cmd
	if !config.NeedsRoot || osinfo.IsRoot() {
		cmd = execCommandInternal(name, arg...)
	} else {
		usesWrapperForPrivileges = true
		if osinfo.CommandExists("sudo") {
			newArgs := []string{"-E"}
			newArgs = append(newArgs, name)
			newArgs = append(newArgs, arg...)
			cmd = execCommandInternal("sudo", newArgs...)
		} else if osinfo.CommandExists("pkexec") {
			newArgs := []string{"--keep-cwd"}
			newArgs = append(newArgs, name)
			newArgs = append(newArgs, arg...)
			cmd = execCommandInternal("pkexec", newArgs...)
		} else {
			return nil, fmt.Errorf("Unable to find app to run as root (sudo or pkexec)")
		}
	}
	

	cmd.Env = append(os.Environ(), config.AdditionalEnv...)

	// Create pipes
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("ExecCommand error stdoutPipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("ExecCommand error stderrPipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("ExecCommand error start: %w", err)
	}

	// Configure stdout writer(s)
	var stdoutWriters []io.Writer
	if config.RetransmitStdout {
		stdoutWriters = append(stdoutWriters, os.Stdout)
	}
	if config.BufferStdout {
		stdoutWriters = append(stdoutWriters, &stdoutBuf)
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
		stderrWriters = append(stderrWriters, &stdoutBuf) // reuse buffer if needed
		if usesWrapperForPrivileges {
			// WARNING
			fmt.Println("BufferStderr: Privilege evaluation was used. The sudo/pkexec prompt will be buffered too")
		}
	}
	if len(stderrWriters) == 0 {
		stderrWriters = []io.Writer{io.Discard}
	}

	// Stream stdout
	go func() {
		io.Copy(io.MultiWriter(stdoutWriters...), stdoutPipe)
	}()

	// Stream stderr
	go func() {
		io.Copy(io.MultiWriter(stderrWriters...), stderrPipe)
	}()

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		if config.ThrowExitCodeError {
			return nil, fmt.Errorf("ExecCommand error command: %w", err)
		}
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return &stdoutBuf, nil
		}
		return nil, fmt.Errorf("ExecCommand error command: %w", err)
	}

	return &stdoutBuf, nil
}
