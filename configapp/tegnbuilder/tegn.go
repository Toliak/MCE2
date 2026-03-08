package tegnbuilder

// Tegn represents an atomic functionality unit, such as a package to install
// or an application to configure (e.g., zsh). It implements TegnGeneral and
// adds parameter configuration capabilities.
type Tegn interface {
	TegnGeneral

	// Returns the list of configurable parameters.
	GetParameters() []TegnParameter

	SetParameter(name string, value string) error
}

type TegnBuildFunc func(data TegnBuilderData) Tegn
