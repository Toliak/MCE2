package main

import (
	"flag"
	"fmt"
	// "os"

	"github.com/toliak/mce/inspector"
)

type Args struct {
	InspectorConfig   inspector.InspectAndHarvestConfig

	InstallDir        string
	Template          string
	Verbosity         int
}

func ParseArgs(args []string) (*Args, error) {
	availableTemplates := []string{"basic", "advanced", "custom"}
	// TODO: where to harvest availableTemplates?
	if len(availableTemplates) == 0 {
		return nil, fmt.Errorf("No available templates")
	}

	checkEnable := flag.Bool("check-enable", true, "Enable platform check")
	harvestEnable := flag.Bool("harvest-enable", true, "Enable harvesting the OS information")
	pkgManagerUpdateEnable := flag.Bool("repo-update-enable", true, "Update the package manager repositories metadata (may require privilege evaluation)")
	pkgManagerGetAvailablePackagesEnable := flag.Bool("repo-packages-enable", true, "Obtain all the available packages from the package manager")

	installDir := flag.String("install-dir", "", "Installation directory")

	var template *string = nil
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

	var verbosity *int = nil
	flag.Func("verbosity", "Verbosity level (0-100, default: 80)", func(s string) error {
		_, err := fmt.Sscanf(s, "%d", verbosity)
		if err != nil {
			return fmt.Errorf("invalid verbosity value: must be an integer")
		}
		if *verbosity < 0 || *verbosity > 100 {
			return fmt.Errorf("verbosity must be between 0 and 100")
		}
		return nil
	})

	// flag.Parse()
	flag.CommandLine.Parse(args)

	// Set defaults
	if template == nil {
		template = &availableTemplates[0]
	}
	if verbosity == nil {
		_verbosity := 80
		verbosity = &_verbosity
	}

	// Return parsed arguments
	return &Args{
		InspectorConfig: inspector.InspectAndHarvestConfig{
			Check: *checkEnable,
			Harvest: *harvestEnable,
			PkgManagerUpdate: *pkgManagerUpdateEnable,
			PkgManagerGetAvailablePackages: *pkgManagerGetAvailablePackagesEnable,
		},
		InstallDir:   *installDir,
		Template:     *template,
		Verbosity:    *verbosity,
	}, nil
}
