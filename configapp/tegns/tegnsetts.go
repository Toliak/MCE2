package tegns

import (
	"github.com/toliak/mce/tegnbuilder"
	"github.com/toliak/mce/tegns/tegnsett"
)

var Tegnsetts = []tegnbuilder.TegnsettBuildFunc{
	tegnsett.NewOSPackages(AllPkgConstructors),
	tegnsett.NewGeneralTegnsett(
		"zsh-config",
		"zsh-config",
		"zsh-config",
		[]string{"os-packages"},
		AllZshConfig,
		nil,
	),
}
