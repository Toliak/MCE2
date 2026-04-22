package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/toliak/mce/inspector"
	"github.com/toliak/mce/platform"
)

type ArgsShared struct {
	Verbosity             int

	InspectorConfig       inspector.InspectAndHarvestConfig
	MainInstallDir        string
	DataDir               string
	UserHomeDir           string

	NoUI                  bool
	ForceConfirm          bool

	SelectEverything      bool
}

type ArgsInstall struct {
	ArgsShared

	MceRepositoryURL      string
	MceRepositoryBranch   string

	JSONPreset			  JSONPreset
}

type ArgsUninstall struct {
	ArgsShared
	JSONPreset			  JSONUninstallPreset
}

func parseSharedArgs(object *ArgsShared) {
	// Inspector flags
	flag.BoolVar(&object.InspectorConfig.Check, "check-enable", true, "Enable platform check")
	flag.BoolVar(&object.InspectorConfig.Harvest, "harvest-enable", true, "Enable harvesting the OS information")
	flag.BoolVar(
		&object.InspectorConfig.PkgManagerUpdate, 
		"repo-update-enable",
		true,
		"Update the package manager repositories metadata (may require privilege evaluation)",
	)
	flag.BoolVar(
		&object.InspectorConfig.PkgManagerGetAvailablePackages,
		"repo-packages-enable",
		true,
		"Obtain all the available packages from the package manager",
	)

	flag.BoolVar(&object.NoUI, "no-ui", false, "Disable UI")
	flag.BoolVar(&object.ForceConfirm, "y", false, "Forcefully confirm installation")
	flag.BoolVar(&object.SelectEverything, "ALL", false, "Select everything in the UI")

	flag.Func("verbosity", "Verbosity level (0-100, default: 80)", func(s string) error {
		var v int
		_, err := fmt.Sscanf(s, "%d", &v)
		if err != nil {
			return fmt.Errorf("invalid verbosity value: must be an integer")
		}
		if v < 0 || v > 100 {
			return fmt.Errorf("verbosity must be between 0 and 100")
		}
		object.Verbosity = v
		return nil
	})

	flag.StringVar(
		&object.UserHomeDir,
		"home-dir",
		object.UserHomeDir,
		"User's home directory",
	)
	flag.StringVar(
		&object.MainInstallDir,
		"main-install-dir",
		filepath.Join("~", ".local", "share", "MakeConfigurationEasier2"),
		"Directory to clone the MCE2 project (absolute path)",
	)
	flag.StringVar(
		&object.DataDir, 
		"data-dir",
		"data", 
		"Path inside '-main-install-dir' where configs and other data will be put",
	)
}

func postParseSharedArgs(object *ArgsShared) error {
	if object.UserHomeDir == "" {
		return fmt.Errorf("Unable to determine the user's home directory. Provide it using '-home-dir' argument")
	}

	splitPath := platform.SplitPath(filepath.Clean(object.MainInstallDir))
	for i, v := range splitPath {
		if v == "~" {
			splitPath[i] = object.UserHomeDir
		}
	}

	return nil
}

func parseInstallArgs(object *ArgsInstall) {
	

	// JSON preset flag
	flag.Func("preset", "JSON preset", func(s string) error {
		preset, err := UnmarshalJSONPreset(s)
		if err != nil {
			return err
		}

		object.JSONPreset = preset
		return nil
	})

	flag.StringVar(&object.MceRepositoryURL, "mce-repo-url", "https://github.com/Toliak/mce2config", "MCE2 repository URL")
	flag.StringVar(&object.MceRepositoryBranch, "mce-repo-branch", "master", "MCE2 branch")
}

func parseUninstallArgs(object *ArgsUninstall) {
	flag.Func("preset", "JSON preset", func(s string) error {
		preset, err := UnmarshalJSONUninstallPreset(s)
		if err != nil {
			return err
		}

		object.JSONPreset = preset
		return nil
	})
}

func ParseInstallArgs(args []string) (*ArgsInstall, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("UserHomeDir error, skipped: %s", err)
		homeDir = ""
	}

	object := ArgsInstall{
		ArgsShared: ArgsShared{
			Verbosity: 80,
			UserHomeDir: homeDir,
		},
		JSONPreset: GetDefaultPreset(),
	}

	parseInstallArgs(&object)
	flag.CommandLine.Parse(args)
	err = postParseSharedArgs(&object.ArgsShared)
	if err != nil {
		return nil, err
	}

	return &object, nil
}

func ParseUninstallArgs(args []string) (*ArgsUninstall, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("UserHomeDir error, skipped: %s", err)
		homeDir = ""
	}

	object := ArgsUninstall{
		ArgsShared: ArgsShared{
			Verbosity: 80,
			UserHomeDir: homeDir,
		},
		JSONPreset: make(JSONUninstallPreset),
	}

	parseUninstallArgs(&object)
	flag.CommandLine.Parse(args)
	err = postParseSharedArgs(&object.ArgsShared)
	if err != nil {
		return nil, err
	}

	return &object, nil
}
