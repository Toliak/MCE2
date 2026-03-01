package main

import (
	"flag"
	"fmt"
	// "os"
)

type Args struct {
	CheckOnly    bool
	CheckDisable bool
	HarvestOnly  bool
	InstallDir   string
	Template     string
	Verbosity    int
}

func ParseArgs(args []string) (*Args, error) {
	availableTemplates := []string{"basic", "advanced", "custom"}
	// TODO: where to harvest availableTemplates?
	if len(availableTemplates) == 0 {
		return nil, fmt.Errorf("No available templates")
	}

	checkOnly := flag.Bool("check-only", false, "Run check only")
	checkDisable := flag.Bool("check-disable", false, "Disable check")
	harvestOnly := flag.Bool("harvest-only", false, "Run harvest only")
	installDir := flag.String("install-dir", "", "Installation directory")

	var template *string = nil
	flag.Func(
		"template", 
		fmt.Sprintf("Package template (one of: %v)", availableTemplates), 
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
		CheckOnly:    *checkOnly,
		CheckDisable: *checkDisable,
		HarvestOnly:  *harvestOnly,
		InstallDir:   *installDir,
		Template:     *template,
		Verbosity:    *verbosity,
	}, nil
}
