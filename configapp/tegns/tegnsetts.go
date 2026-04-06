package tegns

import (
	"github.com/toliak/mce/tegnbuilder"
	"github.com/toliak/mce/tegns/tegnsett"
)

var Tegnsetts = []tegnbuilder.TegnsettBuildFunc{
	tegnsett.NewOSPackages(AllPkgConstructors),
	tegnsett.NewGeneralTegnsett(
		"mce2",
		"MCE2 General",
		"General Make Configuration Easier 2 installation",
		[]string{"os-packages"},
		AllMCE2Tegns,
		nil,
	),
	tegnsett.NewGeneralTegnsett(
		"zsh-config",
		"zsh-config",
		"ZSH configuration",
		[]string{"os-packages", "mce2"},
		AllZshConfig,
		nil,
	),
}
