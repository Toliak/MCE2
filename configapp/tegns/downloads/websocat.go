package downloads

import (
	"fmt"
	"strings"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/tegnbuilder"
	"github.com/toliak/mce/tegns/tegn"
)

var Websocat = tegn.NewGenericDownloadBuilder(
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
