package downloads

import (
	"fmt"
	"strings"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/tegnbuilder"
	"github.com/toliak/mce/tegns/tegn"
)


var Fzf = tegn.NewGenericDownloadBuilder(
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