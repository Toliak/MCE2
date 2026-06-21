package tegns

import (
	"github.com/toliak/mce/tegnbuilder"
	"github.com/toliak/mce/tegns/tegn"
)

// Register packages here
var AllPkgTegns = []tegnbuilder.TegnBuildFunc{
	tegn.NewGenericPackageBuilder(
		"vim",
		"vim",
		"vim",
		"vim",
		nil,
	),
	tegn.NewGenericPackageBuilder(
		"zsh",
		"zsh",
		"zsh",
		"zsh",
		nil,
	),
	tegn.NewGenericPackageBuilder(
		"tmux",
		"tmux",
		"tmux",
		"tmux",
		nil,
	),
	tegn.NewGenericPackageBuilder(
		"bash",
		"bash",
		"bash",
		"bash",
		nil,
	),
	tegn.NewGenericPackageBuilder(
		"git",
		"git",
		"git",
		"git",
		nil,
	),
	tegn.NewGenericPackageBuilder(
		"curl",
		"curl",
		"curl",
		"curl",
		nil,
	),
	tegn.NewGenericPackageBuilder(
		"wget",
		"wget",
		"wget",
		"wget",
		nil,
	),
	tegn.NewGenericPackageBuilder(
		"xclip",
		"xclip",
		"xclip",
		"xclip",
		nil,
	),
	tegn.NewGenericPackageBuilder(
		"net-tools",
		"net-tools",
		"net-tools",
		"net-tools",
		nil,
	),
	tegn.NewGenericPackageBuilder(
		"dnsutils",
		"dnsutils",
		"dnsutils",
		"dnsutils",
		nil,
	),
	tegn.NewGenericPackageBuilder(
		"procps",
		"procps",
		"procps",
		"procps",
		nil,
	),
	tegn.NewGenericPackageBuilder(
		"psmisc",
		"psmisc",
		"psmisc",
		"psmisc",
		nil,
	),
	tegn.NewGenericPackageBuilder(
		"mc",
		"mc",
		"mc",
		"mc",
		nil,
	),
}

// # https://github.com/mjakob-gh/build-static-tmux/releases/tag/v3.5d
// # https://github.com/dtschan/vim-static/blob/master/build.sh
// # https://github.com/romkatv/zsh-bin
// # TODO: move the version into the parameters

var AllDownloadTegns = []tegnbuilder.TegnBuildFunc{
	downloadLf,
	downloadFzf,
	downloadLsGo,
	downloadWebsocat,
}

// Register packages here
var AllZshConfigTegns = []tegnbuilder.TegnBuildFunc{
	tegn.NewTegnZshBaseConfigBuilder(),
	tegn.NewTegnZshPowerLevel10kBuilder(),
	tegn.NewTegnZshSyntaxHighlightBuilder(),
	tegn.NewTegnZshLocalConfigBuilder(),
	tegn.NewTegnZshChshBuilder(),

	tegn.NewTegnZshAutoSuggestionsBuilder(),
}

// Register packages here
var AllBashConfigTegns = []tegnbuilder.TegnBuildFunc{
	tegn.NewTegnBashLocalConfigBuilder(),
}

// Register packages here
var AllMCE2Tegns = []tegnbuilder.TegnBuildFunc{
	tegn.NewTegnCloneRepoBuilder(),
}

// Register packages here
var AllSharedConfigTegns = []tegnbuilder.TegnBuildFunc{
	tegn.NewTegnSharedLocalConfigBuilder(),
}

// Register packages here
var AllVimConfigTegns = []tegnbuilder.TegnBuildFunc{
	tegn.NewTegnUltimateVimBuilder(),
	tegn.NewTegnVimLocalConfigBuilder(),
}

// Register packages here
var AllTmuxConfigTegns = []tegnbuilder.TegnBuildFunc{
	tegn.NewTegnOhMyTmuxBuilder(),
	tegn.NewTegnTmuxLocalConfigBuilder(),
}
