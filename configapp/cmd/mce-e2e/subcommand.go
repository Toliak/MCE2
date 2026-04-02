package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
)

type SubCommandFunc		func (info data.OSInfo, data any) error

// Returns the flagset and the data-object with the parameter pointers
type SubCommandFlagSet	func() (flag.FlagSet, any)

type SubCommandData struct {
	Name			string
	Executor		SubCommandFunc
	FlagSet			SubCommandFlagSet
}

var SubCommands = []SubCommandData {
	{
		Name: "harvest",
		Executor: subCommandHarvestJson,
		FlagSet: emptyFlagSet,
	},
	{
		Name: "repo-list",
		Executor: subCommandRepoList,
		FlagSet: emptyFlagSet,
	},
	{
		Name: "repo-install",
		Executor: subCommandRepoInstall,
		FlagSet: subCommandRepoInstallFlagSet,
	},
	// "harvest": subCommandHarvestJson,
	// "repo-update": subCommandHarvestJson,
	// "repo-search": subCommandHarvestJson,
	// "repo-install": subCommandHarvestJson,
}

func emptyFlagSet() (flag.FlagSet, any) {
	flagSet := flag.NewFlagSet("empty", flag.ExitOnError)
	return *flagSet, struct{}{}
}

func subCommandHarvestJson(info data.OSInfo, _ any) error {
	info_json, err := json.Marshal(info)

	if err != nil {
		return err
	}

	fmt.Printf("%s", info_json)
	return nil
}

func subCommandRepoList(info data.OSInfo, _ any) error {
	err := platform.UpdateRepositories(info.PkgManager)
	if err != nil {
		return fmt.Errorf("Update error: %w", err)
	}

	packageList, err := platform.GetAvailablePackages(
		&info.PkgManager,
	)
	if err != nil {
		return fmt.Errorf("Search error: %w", err)
	}
	info_json, err := json.Marshal(packageList)
	if err != nil {
		return fmt.Errorf("json.Marshal error: %w", err)
	}

	fmt.Printf("%s", info_json)
	return nil
}

type SubCommandRepoInstallData struct {
	// External pointer -- because we need the field to be inside the v.Func in the subCommandRepoInstallFlagSet
	// Internal pointer -- to handle optionality
	toInstall **[]string
}

func subCommandRepoInstallFlagSet() (flag.FlagSet, any) {
	v := flag.NewFlagSet("repo-install", flag.ExitOnError)

	var toInstall *[]string = nil
	v.Func(
		"to-install", 
		"Packages to search and install (split by a colon)", 
		func(s string) error {
			v := strings.Split(s, ":")
			toInstall = &v
			return nil
		},
	)

	return *v, SubCommandRepoInstallData{
		toInstall: &toInstall,
	}
}

func subCommandRepoInstall(info data.OSInfo, data any) error {
	d := data.(SubCommandRepoInstallData)
	if *d.toInstall == nil {
		return fmt.Errorf("Parameter -to-install must be set")
	}

	toInstall := **d.toInstall

	err := platform.UpdateRepositories(info.PkgManager)
	if err != nil {
		return fmt.Errorf("Update error: %w", err)
	}

	found, notFound, err := platform.SearchPackageFullNames(
		info.PkgManager,
		toInstall,
	)
	if err != nil {
		return fmt.Errorf("Search error: %w", err)
	}

	if len(found) != 0 {
		err = platform.InstallPackages(info.PkgManager, found)
		if err != nil {
			return fmt.Errorf("Install error: %w", err)
		}
	}

	info_json, err := json.Marshal(
		map[string][]string{
			"found_and_installed": found,
			"not_found": notFound,
		},
	)
	if err != nil {
		return fmt.Errorf("json.Marshal error: %w", err)
	}

	fmt.Printf("%s", info_json)
	return nil
}
