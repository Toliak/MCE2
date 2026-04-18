package tegn

import (
	// "fmt"

	// "github.com/toliak/mce/inspector"
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/toliak/mce/osinfo/data"
	"github.com/toliak/mce/platform"
	tb "github.com/toliak/mce/tegnbuilder"
)

// type GenericDownloadTegn interface {
// 	tb.Tegn

// 	// Gets the package name
// 	GetPackageName() string
// }

// Raises the error if the URL not available => the Tegn is not available
type GenericDownloadUrlFun func (osInfo tb.OSInfoExt) (string, error)
type GenericDownloadPostProcessFun func (data GenericDownloadPostProcessData) error

type GenericDownloadPostProcessData struct {
	osInfo tb.OSInfoExt
	already tb.TegnInstalledFeaturesMap
	downloadFilePath string
	appDir string
	appName string
}

// The type describes the package that can be installed via the OS package manager
type GenericDownload struct {
	name string
	description string
	getUrl GenericDownloadUrlFun
	postProcess GenericDownloadPostProcessFun
	appName string
	downloadFileName string
}

var _ tb.Tegn = (*GenericDownload)(nil)
// var _ GenericDownloadTegn = (*GenericDownload)(nil)

func NewGenericDownloadBuilder(
	appName string,
	name string,
	description string,
	downloadFileName string,
	getUrl GenericDownloadUrlFun,
	postProcess GenericDownloadPostProcessFun,
) tb.TegnBuildFunc {
	return func () tb.Tegn {
		return &GenericDownload {
			name: name,
			description: description,
			getUrl: getUrl,
			postProcess: postProcess,
			appName: appName,
			downloadFileName: downloadFileName,
		}
	}
}

func getDownloadTempDir() string {
	tempDir := os.TempDir()
	path := filepath.Join(tempDir, "mce2-download-temp")
	return path
}

func getDownloadAppDir() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(userHomeDir, ".local", "bin")
	return path, nil
}

// func (p *GenericDownload) GetPackageName() string {
// 	return p.pkgName
// }

// var _ tb.TegnBuildFunc = NewTegnLinuxPackages
// GetID implements [tegnbuilder.Tegn].
func (p *GenericDownload) GetID() string {
	return "download-" + p.appName
}

// GetName implements [tegnbuilder.Tegn].
func (p *GenericDownload) GetName() string {
	return p.name
}

// GetDescription implements [tegnbuilder.Tegn].
func (p *GenericDownload) GetDescription() string {
	return p.description
}

// GetAvailableCPUArch implements [tegnbuilder.Tegn].
func (p *GenericDownload) GetAvailableCPUArch() *[]data.CPUArchE {
	return nil
}

// GetAvailableOsType implements [tegnbuilder.Tegn].
func (p *GenericDownload) GetAvailableOsType() *[]data.OSTypeE {
	return nil
}

// GetAvailability implements [tegnbuilder.Tegn].
func (p *GenericDownload) GetAvailability(
	osInfo tb.OSInfoExt,
	_before tb.TegnInstalledFeaturesMap,
	enabledIds tb.TegnGeneralEnabledIDsMap,
) tb.TegnAvailability {
	_, err := p.getUrl(osInfo)
	if err != nil {
		return tb.NewTegnNotAvailable(
			fmt.Sprintf("Unable to get URL: %s", err),
		)
	}

	return tb.NewTegnAvailable()
}

// GetFeatures implements [tegnbuilder.Tegn].
func (p *GenericDownload) GetFeatures() tb.TegnInstalledFeaturesMap {
	return tb.TegnInstalledFeaturesMap{
		// TODO: pkg?
		tb.TegnFeature("app:" + p.appName): true,
		tb.TegnFeature("download:" + p.appName): true,
	}
}

// GetBeforeIDs implements [tegnbuilder.Tegn].
func (p *GenericDownload) GetBeforeIDs() []string {
	return make([]string, 0)
}

// GetParameters implements [tegnbuilder.Tegn].
func (p *GenericDownload) GetParameters(osInfo tb.OSInfoExt) []tb.TegnParameter {
	appDir, _ := getDownloadAppDir()

	return[]tb.TegnParameter {
		tb.NewTegnParameter(
			"temp-dir",
			"Temp dir",
			tb.TegnParameterTypeString,
			tb.WithDescription("Download temp directory (read-only)"),
			tb.WithDefaultValue(getDownloadTempDir()),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),
		tb.NewTegnParameter(
			"app-dir",
			"App dir",
			tb.TegnParameterTypeString,
			tb.WithDescription("Target app directory (read-only)"),
			tb.WithDefaultValue(appDir),
			tb.WithAvailabilityFalse("Read-only"),
			tb.WithReadOnlyValidator(),
		),
	}
}

func (p *GenericDownload) IsInstalled(osInfo tb.OSInfoExt) bool {
	downloadDir, err := getDownloadAppDir()
	if err != nil {
		return false
	}

	return platform.FileEntryExists(
		filepath.Join(downloadDir, p.appName),
	)
}

func GenericDownloadPostTarGzUnpack(pathInArchive string) GenericDownloadPostProcessFun {
	return func (data GenericDownloadPostProcessData) error {
		file, err := os.Open(data.downloadFilePath)
		if err != nil {
			return err
		}
		defer file.Close()

		gzr, err := gzip.NewReader(file)
		if err != nil {
			return err
		}
		defer gzr.Close()

		tr := tar.NewReader(gzr)
		
		found := false
		for {
			header, err := tr.Next()

			switch {
			case err == io.EOF:
				return nil
			case err != nil:
				return err
			case header == nil:
				continue
			}

			if header.Name != pathInArchive {
				continue
			}
			found = true
			target := filepath.Join(data.appDir, data.appName)

			switch header.Typeflag {

			case tar.TypeDir:
				continue

			case tar.TypeReg:
				if err := MkdirAllParent(target); err != nil {
					return err
				}

				outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
				if err != nil {
					return err
				}
				defer outFile.Close()

				if _, err := io.Copy(outFile, tr); err != nil {
					return err
				}
			}

			if found {
				break
			}
		}

		if !found {
			return fmt.Errorf("Unable to find file '%s' in the archive '%s'", pathInArchive, data.downloadFilePath)
		}

		return nil
	}
}

func GenericDownloadPostMove(data GenericDownloadPostProcessData) error {
	sourcePath := data.downloadFilePath
	targetPath := filepath.Join(data.appDir, data.appName)
	
	err := os.Rename(sourcePath, targetPath)
	if err != nil {
		return fmt.Errorf("GenericDownloadPostMove Rename: %w", err)
	}
	err = os.Chmod(targetPath, 0755)
	if err != nil {
		return fmt.Errorf("GenericDownloadPostMove Chmod: %w", err)
	}

	return nil
}

func (p *GenericDownload) ExecInstall(osInfo tb.OSInfoExt, already tb.TegnInstalledFeaturesMap, _params tb.TegnParameterMap) error {
	appDir, err := getDownloadAppDir()
	if err != nil {
		return fmt.Errorf("getDownloadAppDir error: %w", err)
	}
	err = os.MkdirAll(appDir, 0755)
	if err != nil {
		return fmt.Errorf("MkdirAll %s: %w", appDir, err)
	}

	tempDownloadDir := getDownloadTempDir()
	err = os.MkdirAll(tempDownloadDir, 0755)
	if err != nil {
		return fmt.Errorf("MkdirAll %s: %w", tempDownloadDir, err)
	}

	url, err := p.getUrl(osInfo)
	if err != nil {
		return fmt.Errorf("getUrl error: %w", err)
	}

	downloadPath := filepath.Join(tempDownloadDir, p.downloadFileName)
	err = platform.DownloadFile(url, downloadPath)
	if err != nil {
		return fmt.Errorf("cannot download '%s' -> '%s': %w", url, downloadPath, err)
	}

	err = p.postProcess(GenericDownloadPostProcessData{
		osInfo: osInfo,
		already: already,
		downloadFilePath: downloadPath,
		appDir: appDir,
		appName: p.appName,
	})
	if err != nil {
		return fmt.Errorf("postProcess: %w", err)
	}

	return nil
}
