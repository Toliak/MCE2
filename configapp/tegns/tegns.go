package tegns

import (
	"github.com/toliak/mce/tegnbuilder"
	"github.com/toliak/mce/tegns/tegn"
)

// Register packages here
var AllPkgConstructors = []tegnbuilder.TegnBuildFunc{
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
}

// Register packages here
var AllZshConfig = []tegnbuilder.TegnBuildFunc{
	tegn.NewTegnZshBaseConfigBuilder(),
	tegn.NewTegnZshPowerLevel10kBuilder(),
}

// Register packages here
var AllMCE2Tegns = []tegnbuilder.TegnBuildFunc{
	tegn.NewTegnCloneRepoBuilder(),
}

// We need categories here
// In the categories there will be the packages

// Register packages here
// var allPkgConstructors = []func(*data.OsInfo) pkg.Pkg{
// 	obj.NewPkgLinuxPackages,
// }

// var allCategories = []func(*data.OsInfo) pkg.Category {
// 	NewGenericCategoryWrap(
// 		GenericCategory{
// 			ID: "package-shit",
// 			Name: "package-shit",
// 			Description: "package-shit",
// 			OSTypes: nil,
// 			CPUArch: nil,
// 			BeforeIDs: make([]string, 0),
// 		},
// 		allPkgConstructors...,
// 	),
// }

// func GetAllCategories(v *data.OsInfo) []pkg.Category {
// 	result := make([]pkg.Category, len(allPkgConstructors))
// 	for i, constructor := range allCategories {
// 		result[i] = constructor(v)
// 	}

// 	return result
// }

// func GetAllPkg(v *data.OsInfo) []pkg.Pkg {
// 	result := make([]pkg.Pkg, len(allPkgConstructors))
// 	for i, constructor := range allPkgConstructors {
// 		result[i] = constructor(v)
// 	}

// 	return result
// }
