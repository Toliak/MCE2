package platform

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/toliak/mce/osinfo/data"
)

// func InstallPackage(info *data.PkgManager, packageName string) error {

// }

// // Return package names
// func SearchPackageNameStartsWith(info *data.PkgManager, packageName string) ([]string, error) {

// }

type PkgManagerError struct {
	what string
}

// Assertion
var _ error = (*PkgManagerError)(nil);

func NewPkgManagerError(what string) *PkgManagerError {
	return &PkgManagerError{
		what: what,
	}
}

func (this* PkgManagerError) Error() string {
	return this.what
}

func execCommandSearch(name string, arg ...string) (*bytes.Buffer, error) {
	v := NewExecCommandWrapper(
		WithBufferStdout(true),
	)
	return ExecCommand(v, name, arg...)
}

func execCommandUpdate(name string, arg ...string) error {
	v := NewExecCommandWrapper(
		WithThrowExitCodeError(true),
	)
	_, err := ExecCommand(v, name, arg...)
	return err
}

func execCommandInstall(name string, arg ...string) error {
	v := NewExecCommandWrapper(
		WithThrowExitCodeError(true),
	)
	_, err := ExecCommand(v, name, arg...)
	return err
}

func UpdateRepositories(info *data.PkgManager) error {
	if info == nil {
		return fmt.Errorf("pkgManager info cannot be nil")
	}

	switch info.V {
	case data.PkgMgrAptGet:
		err := execCommandUpdate(info.Raw, "update", "-y", "--allow-releaseinfo-change")
		return err
	case data.PkgMgrApk:
		err := execCommandUpdate(info.Raw, "update")
		return err
	case data.PkgMgrDnf:
		fmt.Printf("Skipped update for %v package manager\n", info.Raw)
		return nil
	case data.PkgMgrMicroDnf:
		fmt.Printf("Skipped update for %v package manager\n", info.Raw)
		return nil
	case data.PkgMgrYum:
		fmt.Printf("Skipped update for %v package manager\n", info.Raw)
		return nil
	case data.PkgMgrPacman:
		err := execCommandUpdate(info.Raw, "-Syy")
		return err
	case data.PkgMgrBrew:
		return NewPkgManagerError("Not implemented for PkgMgrBrew")
	case data.PkgMgrWinget:
		return NewPkgManagerError("Not implemented for PkgMgrWinget")
	case data.PkgMgrScoop:
		return NewPkgManagerError("Not implemented for PkgMgrScoop")
	default:
		return fmt.Errorf("Unknown package manager %#v", info)
	}
}

func InstallPackages(info *data.PkgManager, packageNames []string) error {
	if info == nil {
		return fmt.Errorf("pkgManager info cannot be nil")
	}

	switch info.V {
	case data.PkgMgrAptGet:
		argList := []string{"install", "-y"}
		argList = append(argList, packageNames...)

		v := NewExecCommandWrapper(
			WithThrowExitCodeError(true),
			WithAdditionalEnv("DEBIAN_FRONTEND=noninteractive"),
		)
		_, err := ExecCommand(v, info.Raw, argList...)

		return err
	case data.PkgMgrApk:
		argList := []string{"add"}
		argList = append(argList, packageNames...)
		err := execCommandInstall(info.Raw, argList...)
		return err
	case data.PkgMgrDnf:
		argList := []string{"install", "-y"}
		argList = append(argList, packageNames...)
		err := execCommandInstall(info.Raw, argList...)
		return err
	case data.PkgMgrMicroDnf:
		argList := []string{"install", "-y"}
		argList = append(argList, packageNames...)
		err := execCommandInstall(info.Raw, argList...)
		return err
	case data.PkgMgrYum:
		argList := []string{"install", "-y"}
		argList = append(argList, packageNames...)
		err := execCommandInstall(info.Raw, argList...)
		return err
	case data.PkgMgrPacman:
		argList := []string{"-S", "--noconfirm"}
		argList = append(argList, packageNames...)
		err := execCommandInstall(info.Raw, argList...)
		return err
	case data.PkgMgrBrew:
		return NewPkgManagerError("Not implemented for PkgMgrBrew")
	case data.PkgMgrWinget:
		return NewPkgManagerError("Not implemented for PkgMgrWinget")
	case data.PkgMgrScoop:
		return NewPkgManagerError("Not implemented for PkgMgrScoop")
	default:
		return fmt.Errorf("Unknown package manager %#v", info)
	}
}

// Return package names
func SearchPackageFullName(info *data.PkgManager, packageName string) (bool, error) {
	found, _, err := SearchPackageFullNames(info, []string{packageName})
	if err != nil {
		return false, err
	}

	return len(found) == 1, nil
}

func searchPackagesDnf(info *data.PkgManager, packageNames []string) ([]string, []string, error) {
	resFound := make([]string, 0)
	resNotFound := make([]string, 0)

	commandsArgs := []string {
		"repoquery",
	}
	commandsArgs = append(commandsArgs, packageNames...)

	fmt.Printf("args: %#v\n", commandsArgs)
	stdoutBuf, err := execCommandSearch("dnf", commandsArgs...)
	if err != nil {
		return resFound, resNotFound, fmt.Errorf("searchPackagesMicroDnf error: %w", err)
	}

	pkgRegexp := regexp.MustCompile(`^(.*)-\d+:.*$`)

	receivedPkgNames := make(map[string]struct{})

	scanner := bufio.NewScanner(stdoutBuf)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		matches := pkgRegexp.FindStringSubmatch(line)
		if len(matches) == 2 {
			receivedPkgNames[matches[1]] = struct{}{};
		}
	}

	for _, v := range packageNames {
		_, ok := receivedPkgNames[v]
		if ok {
			resFound = append(resFound, v)
		} else {
			resNotFound = append(resNotFound, v)
		}
	}

	return resFound, resNotFound, nil
}

func searchPackagesMicroDnf(info *data.PkgManager, packageNames []string) ([]string, []string, error) {
	resFound := make([]string, 0)
	resNotFound := make([]string, 0)

	commandsArgs := []string {
		"repoquery",
	}
	commandsArgs = append(commandsArgs, packageNames...)

	fmt.Printf("args: %#v\n", commandsArgs)
	stdoutBuf, err := execCommandSearch("microdnf", commandsArgs...)
	if err != nil {
		return resFound, resNotFound, fmt.Errorf("searchPackagesMicroDnf error: %w", err)
	}

	pkgRegexp := regexp.MustCompile(`^(.+)-[^-]+-[^-]+\.[^.]+\.[^.]+$`)

	receivedPkgNames := make(map[string]struct{})

	scanner := bufio.NewScanner(stdoutBuf)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		matches := pkgRegexp.FindStringSubmatch(line)
		if len(matches) == 2 {
			receivedPkgNames[matches[1]] = struct{}{};
		}
	}

	for _, v := range packageNames {
		_, ok := receivedPkgNames[v]
		if ok {
			resFound = append(resFound, v)
		} else {
			resNotFound = append(resNotFound, v)
		}
	}

	return resFound, resNotFound, nil
}

func searchPackagesApk(info *data.PkgManager, packageNames []string) ([]string, []string, error) {
	resFound := make([]string, 0)
	resNotFound := make([]string, 0)

	commandsArgs := []string {
		"search",
		"-e",
	}
	commandsArgs = append(commandsArgs, packageNames...)

	fmt.Printf("args: %#v\n", commandsArgs)
	stdoutBuf, err := execCommandSearch("apk", commandsArgs...)
	if err != nil {
		return resFound, resNotFound, fmt.Errorf("searchPackagesApk error: %w", err)
	}

	pkgRegexp := regexp.MustCompile(`^(.+)-[^-]+-r\d+$`)

	receivedPkgNames := make(map[string]struct{})

	scanner := bufio.NewScanner(stdoutBuf)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		matches := pkgRegexp.FindStringSubmatch(line)
		if len(matches) == 2 {
			receivedPkgNames[matches[1]] = struct{}{};
		}
	}

	for _, v := range packageNames {
		_, ok := receivedPkgNames[v]
		if ok {
			resFound = append(resFound, v)
		} else {
			resNotFound = append(resNotFound, v)
		}
	}

	return resFound, resNotFound, nil
}

func searchPackagesYum(info *data.PkgManager, packageNames []string) ([]string, []string, error) {
	resFound := make([]string, 0)
	resNotFound := make([]string, 0)

	commandsArgs := []string {
		"list",
	}
	commandsArgs = append(commandsArgs, packageNames...)
	stdoutBuf, err := execCommandSearch("yum", commandsArgs...)
	if err != nil {
		return resFound, resNotFound, fmt.Errorf("searchPackagesYum error: %w", err)
	}

	pkgRegexp := regexp.MustCompile(`^([^. ]+)\.\S+\s+`)

	receivedPkgNames := make(map[string]struct{})

	scanner := bufio.NewScanner(stdoutBuf)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		matches := pkgRegexp.FindStringSubmatch(line)
		if len(matches) == 2 {
			receivedPkgNames[matches[1]] = struct{}{};
		}
	}

	for _, v := range packageNames {
		_, ok := receivedPkgNames[v]
		if ok {
			resFound = append(resFound, v)
		} else {
			resNotFound = append(resNotFound, v)
		}
	}

	return resFound, resNotFound, nil
}

func searchPackagesPacman(info *data.PkgManager, packageNames []string) ([]string, []string, error) {
	resFound := make([]string, 0)
	resNotFound := make([]string, 0)

	commandsArgs := []string {
		"-Si",
	}
	commandsArgs = append(commandsArgs, packageNames...)
	fmt.Printf("args: %#v\n", commandsArgs)
	stdoutBuf, err := execCommandSearch("pacman", commandsArgs...)
	if err != nil {
		return resFound, resNotFound, fmt.Errorf("searchPackagesPacman error: %w", err)
	}

	pkgRegexp := regexp.MustCompile(`^Name\s+:\s+(.+)`)

	receivedPkgNames := make(map[string]struct{})

	scanner := bufio.NewScanner(stdoutBuf)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		matches := pkgRegexp.FindStringSubmatch(line)
		if len(matches) == 2 {
			receivedPkgNames[matches[1]] = struct{}{};
		}
	}

	for _, v := range packageNames {
		_, ok := receivedPkgNames[v]
		if ok {
			resFound = append(resFound, v)
		} else {
			resNotFound = append(resNotFound, v)
		}
	}

	return resFound, resNotFound, nil
}

func searchPackagesAptCache(info *data.PkgManager, packageNames []string) ([]string, []string, error) {
	resFound := make([]string, 0)
	resNotFound := make([]string, 0)

	// apt-cache --names-only search '^(python3|zsh)$'
	commandsArgs := []string {
		"--names-only",
		"search",
	}
	for i,v := range packageNames {
		packageNames[i] = regexp.QuoteMeta(v);
	}
	joined := strings.Join(packageNames, "|")
	commandsArgs = append(commandsArgs, "^(" + joined + ")$")
	
	v := NewExecCommandWrapper(
		WithBufferStdout(true),
		WithAdditionalEnv("DEBIAN_FRONTEND=noninteractive"),
	)
	stdoutBuf, err := ExecCommand(v, "apt-cache", commandsArgs...)

	if err != nil {
		return resFound, resNotFound, fmt.Errorf("searchPackagesAptCache error: %w", err)
	}

	pkgRegexp := regexp.MustCompile(`^([^ ]+)\s+-\s+`)

	receivedPkgNames := make(map[string]struct{})

	scanner := bufio.NewScanner(stdoutBuf)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		matches := pkgRegexp.FindStringSubmatch(line)
		if len(matches) == 2 {
			receivedPkgNames[matches[1]] = struct{}{};
		}
	}

	for _, v := range packageNames {
		_, ok := receivedPkgNames[v]
		if ok {
			resFound = append(resFound, v)
		} else {
			resNotFound = append(resNotFound, v)
		}
	}

	return resFound, resNotFound, nil
}

func SearchPackageFullNames(info *data.PkgManager, packageNames []string) ([]string, []string, error) {
	if info == nil {
		return make([]string, 0), make([]string, 0), fmt.Errorf("pkgManager info cannot be nil")
	}
	if len(packageNames) == 0 {
		return make([]string, 0), make([]string, 0), nil
	}

	// resFound := make([]string, 0)
	// resNotFound := make([]string, 0)

	fmt.Printf("pkg: %#v", info)

	switch info.V {
	case data.PkgMgrAptGet:
		return searchPackagesAptCache(info, packageNames)
	case data.PkgMgrApk:
		return searchPackagesApk(info, packageNames)
	case data.PkgMgrDnf:
		return searchPackagesMicroDnf(info, packageNames)
	case data.PkgMgrMicroDnf:
		return searchPackagesMicroDnf(info, packageNames)
	case data.PkgMgrYum:
		return searchPackagesYum(info, packageNames)
	case data.PkgMgrPacman:
		return searchPackagesPacman(info, packageNames)
	case data.PkgMgrBrew:
		return nil, nil, NewPkgManagerError("SearchPackageFullNames not implemented for PkgMgrBrew")
	case data.PkgMgrWinget:
		return nil, nil, NewPkgManagerError("SearchPackageFullNames not implemented for PkgMgrWinget")
	case data.PkgMgrScoop:
		return nil, nil, NewPkgManagerError("SearchPackageFullNames not implemented for PkgMgrScoop")
	default:
		return make([]string, 0), make([]string, 0), fmt.Errorf("Unknown package manager %#v", info)
	}
}