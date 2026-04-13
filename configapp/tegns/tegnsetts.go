package tegns

import (
	"github.com/toliak/mce/tegnbuilder"
	"github.com/toliak/mce/tegns/tegnsett"
)

var Tegnsetts = []tegnbuilder.TegnsettBuildFunc{
	tegnsett.NewOSPackages(AllPkgTegns),
	tegnsett.NewGeneralTegnsett(
		"apps-download",
		"Downloaded apps",
		"Applications that can be downloaded without package manager (static binaries usually)",
		[]string{"os-packages"},
		AllDownloadTegns,
		nil,
	),
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
		AllZshConfigTegns,
		nil,
	),
	tegnsett.NewGeneralTegnsett(
		"bash-config",
		"bash-config",
		"Bash configuration",
		[]string{"os-packages", "mce2"},
		AllBashConfigTegns,
		nil,
	),
	tegnsett.NewGeneralTegnsett(
		"shared-shell-config",
		"shared-shell-config",
		"Shared Shell configuration",
		[]string{"zsh-config", "bash-config"},
		AllSharedConfigTegns,
		nil,
	),
	tegnsett.NewGeneralTegnsett(
		"vim-config",
		"vim-config",
		"Vim configuration",
		[]string{"mce2"},
		AllVimConfigTegns,
		nil,
	),
	tegnsett.NewGeneralTegnsett(
		"tmux-config",
		"tmux-config",
		"Tmux configuration",
		[]string{"mce2"},
		AllTmuxConfigTegns,
		nil,
	),
}
