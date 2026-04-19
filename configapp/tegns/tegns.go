package tegns

import (
	"fmt"
	"strings"

	"github.com/toliak/mce/osinfo/data"
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
	tegn.NewGenericDownloadBuilder(
		"lf",
		"lf r41",
		"list files GoLang app r41",
		"lf-r41.tar.gz",
		func(osInfo tegnbuilder.OSInfoExt) (string, error) {
			var sb strings.Builder
			// TODO: verify checksums
			sb.WriteString("https://github.com/gokcehan/lf/releases/download/r41/lf-")

			switch osInfo.OsType.V {
			case data.OSTypeLinux: 
				sb.WriteString("linux-")
				switch osInfo.Arch.V {
				case data.CPUArchAMD64:
					sb.WriteString("amd64")
				case data.CPUArchI386:
					sb.WriteString("386")
				case data.CPUArchAARCH64:
					sb.WriteString("arm64")
				case data.CPUArchARMv7:
					sb.WriteString("arm")
				case data.CPUArchMIPS64Le:
					sb.WriteString("mips64le")
				case data.CPUArchPPC64:
					sb.WriteString("ppc64")
				default:
					return "", fmt.Errorf("OSTypeLinux Architecture not available: %s", &osInfo.Arch)
				}
			case data.OSTypeAndroid: 
				sb.WriteString("android-")
				switch osInfo.Arch.V {
				case data.CPUArchAARCH64:
					sb.WriteString("arm64")
				default:
					return "", fmt.Errorf("OSTypeAndroid Architecture not available: %s", &osInfo.Arch)
				}
			// TODO: support windows's zip archive
			// case data.OSTypeWindows: 
			// 	sb.WriteString("windows-")
			// 	switch osInfo.Arch.V {
			// 	case data.CPUArchAMD64:
			// 		sb.WriteString("amd64")
			// 	case data.CPUArchI386:
			// 		sb.WriteString("386")
			// 	default:
			// 		return "", fmt.Errorf("OSTypeWindows Architecture not available: %s", &osInfo.Arch)
			// 	}
			case data.OSTypeDarwin: 
				sb.WriteString("darwin-")
				switch osInfo.Arch.V {
				case data.CPUArchAMD64:
					sb.WriteString("amd64")
				case data.CPUArchAARCH64:
					sb.WriteString("arm64")
				default:
					return "", fmt.Errorf("OSTypeDarwin Architecture not available: %s", &osInfo.Arch)
				}
			default:
				return "", fmt.Errorf("OSType not available: %s", &osInfo.OsType)
			}

			sb.WriteString(".tar.gz")
			return sb.String(), nil
		},
		tegn.GenericDownloadPostTarGzUnpack("lf"),
	),
	tegn.NewGenericDownloadBuilder(
		"fzf",
		"fzf 0.71.0",
		"fzf GoLang app 0.71.0",
		"fzf-0.71.0.tar.gz",
		func(osInfo tegnbuilder.OSInfoExt) (string, error) {
			var sb strings.Builder
			// TODO: verify checksums
			sb.WriteString("https://github.com/junegunn/fzf/releases/download/v0.71.0/fzf-0.71.0-")

			switch osInfo.OsType.V {
			case data.OSTypeLinux: 
				sb.WriteString("linux_")
				switch osInfo.Arch.V {
				case data.CPUArchAMD64:
					sb.WriteString("amd64")
				case data.CPUArchAARCH64:
					sb.WriteString("arm64")
				case data.CPUArchARMv7:
					sb.WriteString("armv7")
				case data.CPUArchRISCV64:
					sb.WriteString("riscv64")
				default:
					return "", fmt.Errorf("OSTypeLinux Architecture not available: %s", &osInfo.Arch)
				}
			case data.OSTypeAndroid: 
				sb.WriteString("android_")
				switch osInfo.Arch.V {
				case data.CPUArchAARCH64:
					sb.WriteString("arm64")
				default:
					return "", fmt.Errorf("OSTypeAndroid Architecture not available: %s", &osInfo.Arch)
				}
			// TODO: support windows's zip archive
			case data.OSTypeDarwin: 
				sb.WriteString("darwin_")
				switch osInfo.Arch.V {
				case data.CPUArchAMD64:
					sb.WriteString("amd64")
				case data.CPUArchAARCH64:
					sb.WriteString("arm64")
				default:
					return "", fmt.Errorf("OSTypeDarwin Architecture not available: %s", &osInfo.Arch)
				}
			default:
				return "", fmt.Errorf("OSType not available: %s", &osInfo.OsType)
			}

			sb.WriteString(".tar.gz")
			return sb.String(), nil
		},
		tegn.GenericDownloadPostTarGzUnpack("fzf"),
	),
	tegn.NewGenericDownloadBuilder(
		"ls-go",
		"ls-go 1.0.2",
		"ls-go GoLang app 1.0.2",
		"ls-go",
		func(osInfo tegnbuilder.OSInfoExt) (string, error) {
			var sb strings.Builder
			// TODO: verify checksums
			sb.WriteString("https://github.com/acarl005/ls-go/releases/download/v1.0.2/ls-go-")

			switch osInfo.OsType.V {
			case data.OSTypeLinux: 
				sb.WriteString("linux-")
				switch osInfo.Arch.V {
				case data.CPUArchAMD64:
					sb.WriteString("amd64")
				case data.CPUArchI386:
					sb.WriteString("386")
				case data.CPUArchAARCH64:
					sb.WriteString("arm64")
				default:
					return "", fmt.Errorf("OSTypeLinux Architecture not available: %s", &osInfo.Arch)
				}
			case data.OSTypeDarwin: 
				sb.WriteString("darwin-")
				switch osInfo.Arch.V {
				case data.CPUArchAMD64:
					sb.WriteString("amd64")
				case data.CPUArchAARCH64:
					sb.WriteString("arm64")
				default:
					return "", fmt.Errorf("OSTypeDarwin Architecture not available: %s", &osInfo.Arch)
				}
			default:
				return "", fmt.Errorf("OSType not available: %s", &osInfo.OsType)
			}

			return sb.String(), nil
		},
		tegn.GenericDownloadPostMove,
	),
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
