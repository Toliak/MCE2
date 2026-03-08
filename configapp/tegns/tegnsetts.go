package tegns

import (
	"github.com/toliak/mce/tegnbuilder"
	"github.com/toliak/mce/tegns/tegnsett"
)

var Tegnsetts = []tegnbuilder.TegnsettBuildFunc {
	tegnsett.NewOuterOSPackages(AllPkgConstructors),
}

func InitializeAllTegnsetts(tegnsetts []tegnbuilder.TegnsettBuildFunc, data tegnbuilder.TegnBuilderData) []tegnbuilder.Tegnsett {
	result := make([]tegnbuilder.Tegnsett, len(tegnsetts))
	for i, v := range tegnsetts {
		result[i] = v(data)
	}

	return result
}
