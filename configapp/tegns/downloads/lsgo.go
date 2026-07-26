package downloads

import (
	"fmt"
	"strings"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/tegnbuilder"
	"github.com/toliak/mce/tegns/tegn"
)

// TODO: ls-go is not statically built
var LsGo = tegn.NewGenericDownloadBuilder(
	"ls-go",
	"ls-go 1.0.2",
	"ls-go GoLang app 1.0.2",
	"ls-go",
	func(osInfo tegnbuilder.OSInfoExt) (string, error) {
		// TODO: return struct with {name, type: {exefile, zip, targz, tarzst, ...}, where-the-executable-if-archive}
		// TODO: FIXME: ^^^^
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
