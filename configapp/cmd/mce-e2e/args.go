package main

// import (
// 	"flag"
// 	"fmt"
// 	// "os"
// )

// func ParseArgs(args []string) (*SubCommandData, error) {
// 	if len(SubCommands) == 0 {
// 		return nil, fmt.Errorf("No available subcommands")
// 	}
// 	subCommandsMap := make(map[string]SubCommandData, len(SubCommands))
// 	availableSubCommands := make([]string, len(SubCommands))
// 	for i, v := range SubCommands {
// 		availableSubCommands[i] = v.Name
// 		subCommandsMap[v.Name] = v
// 	}

// 	if len(args) < 1 {
// 		return nil, fmt.Errorf("Expected one of the subcommands %v", availableSubCommands)
// 	}

// 	subCommand := args[1]
// }
