package tegnbuilder

// Tegnsett represents a group of related Tegns.
type Tegnsett interface {
	TegnGeneral

	// Returns all child Tegns in this group.
	GetChildren() []Tegn
}

type TegnsettBuildFunc func(data TegnBuilderData) Tegnsett
type TegnsettOuterBuildFunc func (children []TegnBuildFunc) TegnsettBuildFunc
