package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/toliak/mce/inspector"
)

type Args struct {
	InspectorConfig       inspector.InspectAndHarvestConfig

	Template              string
	Verbosity             int

	MainInstallDir        string
	DataDir               string
	MceRepositoryURL      string
	MceRepositoryBranch   string

	JSONPreset			  JSONPreset
	NoUI                  bool
	ForceConfirm          bool
}

func ParseArgs(args []string) (*Args, error) {
	availableTemplates := []string{"basic", "advanced", "custom"}
	// TODO: where to harvest availableTemplates?
	if len(availableTemplates) == 0 {
		return nil, fmt.Errorf("no available templates")
	}

	// Inspector flags
	checkEnable := flag.Bool("check-enable", true, "Enable platform check")
	harvestEnable := flag.Bool("harvest-enable", true, "Enable harvesting the OS information")
	pkgManagerUpdateEnable := flag.Bool(
		"repo-update-enable",
		true,
		"Update the package manager repositories metadata (may require privilege evaluation)",
	)
	pkgManagerGetAvailablePackagesEnable := flag.Bool(
		"repo-packages-enable",
		true,
		"Obtain all the available packages from the package manager",
	)

	noUI := flag.Bool("no-ui", false, "Disable UI")
	forceConfirm := flag.Bool("y", false, "Forcefully confirm installation")

	// Template flag
	var template *string
	flag.Func(
		"template",
		fmt.Sprintf("Selected Tegns template (one of: %v)", availableTemplates),
		func(s string) error {
			for _, t := range availableTemplates {
				if s == t {
					template = &s
					return nil
				}
			}
			return fmt.Errorf("template must be one of %v", availableTemplates)
		},
	)

	// Verbosity flag
	var verbosity *int
	flag.Func("verbosity", "Verbosity level (0-100, default: 80)", func(s string) error {
		var v int
		_, err := fmt.Sscanf(s, "%d", &v)
		if err != nil {
			return fmt.Errorf("invalid verbosity value: must be an integer")
		}
		if v < 0 || v > 100 {
			return fmt.Errorf("verbosity must be between 0 and 100")
		}
		verbosity = &v
		return nil
	})

	// Verbosity flag
	var jsonPreset *JSONPreset
	flag.Func("preset", "JSON preset", func(s string) error {
		preset, err := UnmarshalJSONPreset(s)
		if err != nil {
			return err
		}

		jsonPreset = &preset
		return nil
	})
	if jsonPreset == nil {
		preset := GetDefaultPreset()
		jsonPreset = &preset
	}

	// New MCE2 flags
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("ParseArgs: UserHomeDir error: %w", err)
	}
	mainInstallDir := flag.String(
		"main-install-dir",
		filepath.Join(homeDir, ".local", "share", "MakeConfigurationEasier2"),
		"Directory to clone the MCE2 project (absolute path)",
	)
	dataDir := flag.String("data-dir", "data", "Path inside MainInstallDir where configs and other data will be put")
	mceRepoURL := flag.String("mce-repo-url", "https://github.com/Toliak/MCE2", "MCE2 repository URL")
	mceRepoBranch := flag.String("mce-repo-branch", "master", "MCE2 branch")

	// Parse
	flag.CommandLine.Parse(args)

	// Set defaults
	if template == nil {
		template = &availableTemplates[0]
	}
	if verbosity == nil {
		v := 80
		verbosity = &v
	}

	// Return parsed arguments
	return &Args{
		InspectorConfig: inspector.InspectAndHarvestConfig{
			Check:                         *checkEnable,
			Harvest:                       *harvestEnable,
			PkgManagerUpdate:               *pkgManagerUpdateEnable,
			PkgManagerGetAvailablePackages: *pkgManagerGetAvailablePackagesEnable,
		},
		Template:            *template,
		Verbosity:           *verbosity,
		MainInstallDir:      *mainInstallDir,
		DataDir:             *dataDir,
		MceRepositoryURL:    *mceRepoURL,
		MceRepositoryBranch: *mceRepoBranch,
		JSONPreset:          *jsonPreset,
		NoUI: *noUI,
		ForceConfirm: *forceConfirm,
	}, nil
}
