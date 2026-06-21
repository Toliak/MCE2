package tegns

import (
	"fmt"
	"strings"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/tegnbuilder"
	"github.com/toliak/mce/tegns/tegn"
)

var downloadLf = tegn.NewGenericDownloadBuilder(
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
)

var downloadFzf = tegn.NewGenericDownloadBuilder(
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
)

// TODO: ls-go is not statically built
var downloadLsGo = tegn.NewGenericDownloadBuilder(
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
)

var downloadWebsocat = tegn.NewGenericDownloadBuilder(
	"websocat",
	"websocat 1.14.1",
	"Command-line client for WebSockets",
	"websocat",
	func(osInfo tegnbuilder.OSInfoExt) (string, error) {
		var sb strings.Builder
		sb.WriteString("https://github.com/vi/websocat/releases/download/v1.14.1/websocat.")

		switch osInfo.OsType.V {
		case data.OSTypeLinux: 
			switch osInfo.Arch.V {
			case data.CPUArchAMD64:
				sb.WriteString("x86_64-unknown-linux-musl")
			case data.CPUArchAARCH64:
				sb.WriteString("aarch64-unknown-linux-musl")
			case data.CPUArchARMv7:
				sb.WriteString("arm-unknown-linux-musleabi")
			case data.CPUArchRISCV64:
				sb.WriteString("riscv64gc-unknown-linux-musl")
			case data.CPUArchI386:
				sb.WriteString("i686-unknown-linux-musl")
			default:
				return "", fmt.Errorf("OSTypeLinux Architecture not available: %s", &osInfo.Arch)
			}
		case data.OSTypeAndroid: 
			switch osInfo.Arch.V {
			case data.CPUArchAARCH64:
				sb.WriteString("aarch64-linux-android")
			case data.CPUArchARMv7:
				sb.WriteString("armv7-linux-androideabi")
			default:
				return "", fmt.Errorf("OSTypeAndroid Architecture not available: %s", &osInfo.Arch)
			}
		case data.OSTypeDarwin: 
			switch osInfo.Arch.V {
			case data.CPUArchAMD64:
				sb.WriteString("x86_64")
			case data.CPUArchAARCH64:
				sb.WriteString("aarch64")
			default:
				return "", fmt.Errorf("OSTypeDarwin Architecture not available: %s", &osInfo.Arch)
			}
			sb.WriteString("-apple-darwin")
		case data.OSTypeWindows: 
			switch osInfo.Arch.V {
			case data.CPUArchAMD64:
				sb.WriteString("x86_64")
			case data.CPUArchI386:
				sb.WriteString("i686")
			default:
				return "", fmt.Errorf("OSTypeDarwin Architecture not available: %s", &osInfo.Arch)
			}
			sb.WriteString("-pc-windows-gnu.exe")
		default:
			return "", fmt.Errorf("OSType not available: %s", &osInfo.OsType)
		}

		return sb.String(), nil
	},
	tegn.GenericDownloadPostMove,
)
