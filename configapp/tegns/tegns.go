package tegns

import (
	"github.com/toliak/mce/tegnbuilder"
	"github.com/toliak/mce/tegns/tegn"
)

// Register packages here
var AllPkgConstructors = []tegnbuilder.TegnBuildFunc{
	tegn.NewPkgLinuxPackages,
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
