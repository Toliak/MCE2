package tegnbuilder

// Tegn represents an atomic functionality unit, such as a package to install
// or an application to configure (e.g., zsh). It implements TegnGeneral and
// adds parameter configuration capabilities.
type Tegn interface {
	TegnGeneral

	// Returns the list of configurable parameters.
	// osInfo must not be nil.
	GetParameters(osInfo *OSInfoExt) []TegnParameter

	// The features parameter contains features provided 
	// by previously selected Tegns (according to installation order).
	// SetContextFeatures(features []string)

	// Returns the features provided by this Tegn.
	// Features do not depend on the parameters, OSInfo and before-installed featured.
	GetFeatures() TegnInstalledFeaturesMap

	// Returns true if the Tegn already installed.
	// osInfo must not be nil.
	IsInstalled(osInfo *OSInfoExt) bool

	// Install the Tegn.
	// WARN: this function does not check, is Tegn already installed.
	// Params must contain all parameters value (even not available!)
	// osInfo must not be nil.
	ExecInstall(osInfo *OSInfoExt, already TegnInstalledFeaturesMap, params TegnParameterMap) error

	// WARN: this function does not check, is Tegn already installed
	// ExecUpdate() error

	// WARN: this function does not check, is Tegn already installed
	// osInfo must not be nil.
	ExecUninstall(osInfo *OSInfoExt) error
}

type TegnBuildFunc func() Tegn
