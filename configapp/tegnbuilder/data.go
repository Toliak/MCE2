package tegnbuilder

import "github.com/toliak/mce/osinfo/data"

type TegnBuilderData struct {
	data.OSInfo

	AvailableManagerPackages []string
}
