package main

import (
	"fmt"
	"os"

	"github.com/toliak/mce/inspector"
)

func mainInternal() error {
	data, err := inspector.InspectAndHarvest(inspector.InspectAndHarvestConfig{
		Check: true,
		Harvest: true,
		PkgManagerUpdate: false,
		PkgManagerGetAvailablePackages: false,
	})
	if err != nil {
		return err
	}

	if data == nil {
		return fmt.Errorf("No harvest data obtained, internal error")
	}
	if data.OSInfo == nil {
		return fmt.Errorf("No OSInfo data obtained")
	}

	args := os.Args[1:]

	if len(SubCommands) == 0 {
		return fmt.Errorf("No available subcommands")
	}
	subCommandsMap := make(map[string]SubCommandData, len(SubCommands))
	availableSubCommands := make([]string, len(SubCommands))
	for i, v := range SubCommands {
		availableSubCommands[i] = v.Name
		subCommandsMap[v.Name] = v
	}

	if len(args) < 1 {
		return fmt.Errorf("Expected one of the subcommands %v", availableSubCommands)
	}

	subCommand := args[0]
	subCommandData, ok := subCommandsMap[subCommand]
	if !ok {
		return fmt.Errorf("SubCommand '%s' not found. Available: %v", subCommand, availableSubCommands)
	}

	flagSet, flagSetData := subCommandData.FlagSet()
	err = flagSet.Parse(args[1:])
	if err != nil {
		return err
	}

	return subCommandData.Executor(*data.OSInfo, flagSetData)
}

func main() {
	err := mainInternal()
	if err != nil {
		fmt.Printf("ERROR\n")
		fmt.Printf("%s\n", err)
		os.Exit(1)
		return
	}
}
