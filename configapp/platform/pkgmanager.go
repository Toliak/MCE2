package platform

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/toliak/mce/osinfo/data"
)

type PkgManagerError struct {
	what string
}

var _ error = (*PkgManagerError)(nil)


func (e *PkgManagerError) Error() string {
	return e.what
}

func execCommandSearchWrapper() ExecCommandWrapper {
	return NewExecCommandWrapper(
		WithBufferStdout(true),
	)
}

func execCommandUpdate(name string, arg ...string) error {
	v := NewExecCommandWrapper(
		WithThrowExitCodeError(true),
		WithNeedsRoot(true),
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
		return execCommandUpdate(info.Raw, "update", "-y", "--allow-releaseinfo-change")
	case data.PkgMgrApk:
		return execCommandUpdate(info.Raw, "update")
	case data.PkgMgrDnf, data.PkgMgrMicroDnf, data.PkgMgrYum:
		fmt.Printf("Skipped update for %v package manager\n", info.Raw)
		return nil
	case data.PkgMgrPacman:
		return execCommandUpdate(info.Raw, "-Syy")
	case data.PkgMgrBrew, data.PkgMgrWinget, data.PkgMgrScoop:
		return &PkgManagerError{what: fmt.Sprintf("Not implemented for %v", info.V)}
	default:
		return fmt.Errorf("Unknown package manager %#v", info)
	}
}

func InstallPackages(info *data.PkgManager, packageNames []string) error {
	if info == nil {
		return fmt.Errorf("pkgManager info cannot be nil")
	}

	var cmd string
	var argList []string
	var env []string

	switch info.V {
	case data.PkgMgrAptGet:
		cmd = info.Raw
		argList = []string{"install", "-y"}
		env = []string{"DEBIAN_FRONTEND=noninteractive"}
	case data.PkgMgrApk:
		cmd = info.Raw
		argList = []string{"add"}
	case data.PkgMgrDnf, data.PkgMgrMicroDnf, data.PkgMgrYum:
		cmd = info.Raw
		argList = []string{"install", "-y"}
	case data.PkgMgrPacman:
		cmd = info.Raw
		argList = []string{"-S", "--noconfirm"}
	case data.PkgMgrBrew, data.PkgMgrWinget, data.PkgMgrScoop:
		return &PkgManagerError{what: fmt.Sprintf("Not implemented for %v", info.V)}
	default:
		return fmt.Errorf("Unknown package manager %#v", info)
	}

	argList = append(argList, packageNames...)

	v := NewExecCommandWrapper(
		WithThrowExitCodeError(true),
		WithAdditionalEnvList(env),
		WithNeedsRoot(true),
	)
	_, err := ExecCommand(v, cmd, argList...)
	return err
}

// WARNING: do not use this function if you have to search multiple packages
// Use `SearchPackageFullNames` instead
func SearchPackageFullName(info *data.PkgManager, packageName string) (bool, error) {
	found, _, err := SearchPackageFullNames(info, []string{packageName})
	if err != nil {
		return false, err
	}

	return len(found) == 1, nil
}

type searchConfig struct {
	command      string
	baseArgs     []string
	pkgRegex     *regexp.Regexp
	regexGroup   int
	processNames func([]string) []string
	execWrapper  func() ExecCommandWrapper
}

func searchPackagesGeneric(info *data.PkgManager, packageNames []string, config searchConfig) ([]string, []string, error) {
	resFound := make([]string, 0)
	resNotFound := make([]string, 0)

	commandsArgs := make([]string, len(config.baseArgs))
	copy(commandsArgs, config.baseArgs)

	processedNames := packageNames
	if config.processNames != nil {
		processedNames = config.processNames(packageNames)
	}
	commandsArgs = append(commandsArgs, processedNames...)

	// fmt.Printf("args: %#v\n", commandsArgs)

	var stdoutBuf *bytes.Buffer
	var err error

	var execWrapper func() ExecCommandWrapper = nil
	if config.execWrapper != nil {
		execWrapper = config.execWrapper
	} else {
		execWrapper = execCommandSearchWrapper
	}
	
	stdoutBuf, err = ExecCommand(
		execWrapper(), 
		config.command, 
		commandsArgs...,
	)

	if err != nil {
		return resFound, resNotFound, fmt.Errorf("searchPackages (%s) error: %w", config.command, err)
	}

	receivedPkgNames := make(map[string]struct{})
	scanner := bufio.NewScanner(stdoutBuf)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		matches := config.pkgRegex.FindStringSubmatch(line)
		if len(matches) > config.regexGroup {
			receivedPkgNames[matches[config.regexGroup]] = struct{}{}
		}
	}

	for _, v := range packageNames {
		if _, ok := receivedPkgNames[v]; ok {
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

	// fmt.Printf("pkg: %#v\n", info)

	switch info.V {
	case data.PkgMgrAptGet:
		processedNames := make([]string, len(packageNames))
		for i, v := range packageNames {
			processedNames[i] = regexp.QuoteMeta(v)
		}
		joined := strings.Join(processedNames, "|")
		return searchPackagesGeneric(info, packageNames, searchConfig{
			command:    "apt-cache",
			baseArgs:   []string{"--names-only", "search", "^(" + joined + ")$"},
			pkgRegex:   regexp.MustCompile(`^([^ ]+)\s+-\s+`),
			regexGroup: 1,
			execWrapper: func() ExecCommandWrapper {
				return NewExecCommandWrapper(
					WithBufferStdout(true),
					WithAdditionalEnv("DEBIAN_FRONTEND=noninteractive"),
				)
			},
		})
	case data.PkgMgrApk:
		return searchPackagesGeneric(info, packageNames, searchConfig{
			command:    "apk",
			baseArgs:   []string{"search", "-e"},
			pkgRegex:   regexp.MustCompile(`^(.+)-[^-]+-r\d+$`),
			regexGroup: 1,
		})
	case data.PkgMgrDnf:
		return searchPackagesGeneric(info, packageNames, searchConfig{
			command:    "dnf",
			baseArgs:   []string{"repoquery"},
			pkgRegex:   regexp.MustCompile(`^(.*)-\d+:.*$`),
			regexGroup: 1,
		})
	case data.PkgMgrMicroDnf:
		return searchPackagesGeneric(info, packageNames, searchConfig{
			command:    "microdnf",
			baseArgs:   []string{"repoquery"},
			pkgRegex:   regexp.MustCompile(`^(.+)-[^-]+-[^-]+\.[^.]+\.[^.]+$`),
			regexGroup: 1,
		})
	case data.PkgMgrYum:
		return searchPackagesGeneric(info, packageNames, searchConfig{
			command:    "yum",
			baseArgs:   []string{"list"},
			pkgRegex:   regexp.MustCompile(`^([^. ]+)\.\S+\s+`),
			regexGroup: 1,
		})
	case data.PkgMgrPacman:
		return searchPackagesGeneric(info, packageNames, searchConfig{
			command:    "pacman",
			baseArgs:   []string{"-Si"},
			pkgRegex:   regexp.MustCompile(`^Name\s+:\s+(.+)`),
			regexGroup: 1,
		})
	case data.PkgMgrBrew, data.PkgMgrWinget, data.PkgMgrScoop:
		return nil, nil, &PkgManagerError{what: fmt.Sprintf("SearchPackageFullNames not implemented for %v", info.V)}
	default:
		return make([]string, 0), make([]string, 0), &PkgManagerError{what: fmt.Sprintf("Unknown package manager %#v", info)}
	}
}