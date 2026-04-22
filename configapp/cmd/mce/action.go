package main

import (
	"fmt"
	"strings"
)

// ActionType represents the type of action (enum)
type ActionType int

const (
	ActionInstall ActionType = iota // Install = 0
	ActionUninstall                 // Uninstall = 1
)

// String returns the string representation of the ActionType
func (a ActionType) String() string {
	switch a {
	case ActionInstall:
		return "Install"
	case ActionUninstall:
		return "Uninstall"
	default:
		return fmt.Sprintf("ActionType(%d)", int(a))
	}
}

// ParseActionType converts a string to ActionType
func ParseActionType(s string) (ActionType, error) {
	switch strings.ToUpper(s) {
	case "INSTALL":
		return ActionInstall, nil
	case "UNINSTALL":
		return ActionUninstall, nil
	case "REMOVE":
		return ActionUninstall, nil
	default:
		return -1, fmt.Errorf("invalid action type: %s", s)
	}
}

// Returns ActionType, next argv (1 or 2), error
func ParseActionTypeFromArgv1(argv1 string) (ActionType, int, error) {
	if strings.HasPrefix(argv1, "-") {
		// No action specified, install by default
		return ActionInstall, 1, nil
	}

	actionType, err := ParseActionType(argv1)
	if err != nil {
		return -1, 1, err
	}

	return actionType, 2, nil
}