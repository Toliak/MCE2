package tegns

import (
	"github.com/toliak/mce/optional"
	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/tegnbuilder"
	"github.com/toliak/mce/tegns/downloads"
	"github.com/toliak/mce/tegns/tegn"
)

// Register packages here
var AllPkgTegns = []tegnbuilder.TegnBuildFunc{
	tegn.NewGenericPackageBuilder(
		"vim",
		"vim",
		"vim",
		"vim",
		optional.None[[]data.OSTypeE](),
	),
	tegn.NewGenericPackageBuilder(
		"zsh",
		"zsh",
		"zsh",
		"zsh",
		optional.None[[]data.OSTypeE](),
	),
	tegn.NewGenericPackageBuilder(
		"tmux",
		"tmux",
		"tmux",
		"tmux",
		optional.None[[]data.OSTypeE](),
	),
	tegn.NewGenericPackageBuilder(
		"bash",
		"bash",
		"bash",
		"bash",
		optional.None[[]data.OSTypeE](),
	),
	tegn.NewGenericPackageBuilder(
		"git",
		"git",
		"git",
		"git",
		optional.None[[]data.OSTypeE](),
	),
	tegn.NewGenericPackageBuilder(
		"curl",
		"curl",
		"curl",
		"curl",
		optional.None[[]data.OSTypeE](),
	),
	tegn.NewGenericPackageBuilder(
		"wget",
		"wget",
		"wget",
		"wget",
		optional.None[[]data.OSTypeE](),
	),
	tegn.NewGenericPackageBuilder(
		"xclip",
		"xclip",
		"xclip",
		"xclip",
		optional.None[[]data.OSTypeE](),
	),
	tegn.NewGenericPackageBuilder(
		"net-tools",
		"net-tools",
		"net-tools",
		"net-tools",
		optional.None[[]data.OSTypeE](),
	),
	tegn.NewGenericPackageBuilder(
		"dnsutils",
		"dnsutils",
		"dnsutils",
		"dnsutils",
		optional.None[[]data.OSTypeE](),
	),
	tegn.NewGenericPackageBuilder(
		"procps",
		"procps",
		"procps",
		"procps",
		optional.None[[]data.OSTypeE](),
	),
	tegn.NewGenericPackageBuilder(
		"psmisc",
		"psmisc",
		"psmisc",
		"psmisc",
		optional.None[[]data.OSTypeE](),
	),
	tegn.NewGenericPackageBuilder(
		"mc",
		"mc",
		"mc",
		"mc",
		optional.None[[]data.OSTypeE](),
	),
}

// # https://github.com/mjakob-gh/build-static-tmux/releases/tag/v3.5d
// # https://github.com/dtschan/vim-static/blob/master/build.sh
// # https://github.com/romkatv/zsh-bin
// # TODO: move the version into the parameters

var AllDownloadTegns = []tegnbuilder.TegnBuildFunc{
	downloads.Lf,
	downloads.Fzf,
	downloads.LsGo,
	downloads.Websocat,
	downloads.Gitstatusd,
}

// Register packages here
var AllZshConfigTegns = []tegnbuilder.TegnBuildFunc{
	tegn.NewTegnZshBaseConfigBuilder(),
	tegn.NewTegnZshPowerLevel10kBuilder(),
	tegn.NewTegnZshSyntaxHighlightBuilder(),
	tegn.NewTegnZshFavoriteDirBuilder(),
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
