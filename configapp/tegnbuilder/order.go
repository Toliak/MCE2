package tegnbuilder

import (
	"fmt"
	"slices"
)

type OrderError struct {
	what string
}

var _ error = (*OrderError)(nil)

func (e *OrderError) Error() string {
	return e.what
}

func GetPkgGeneralOrder(packages []TegnGeneral) ([]string, error) {
	packageToBeforeIDs := make(map[string][]string, len(packages))

	allBeforeIDs := make([]string, 0, len(packages))
	for _, v := range packages {
		beforeIDs := v.GetBeforeIDs()
		packageToBeforeIDs[v.GetID()] = beforeIDs
		allBeforeIDs = append(allBeforeIDs, beforeIDs...)
	}

	for _, id := range allBeforeIDs {
		if _, ok := packageToBeforeIDs[id]; ok {
			continue
		}

		fmt.Printf("BeforeID '%s' not found, ignored\n", id)

		// Remove that ID from any dependencies
		for k, v := range packageToBeforeIDs {
			index := slices.Index(v, id)
			if index == -1 {
				continue
			}

			packageToBeforeIDs[k] = slices.Delete(v, index, index+1)
		}
	}

	// Find packages (in packageToBeforeIDs) with empty dependencies -- insert them into the result array.
	// Remove that package id from the BeforeID and from the another package dependencies.
	// Do that until packageToBeforeIDs is empty (all packages processed).
	// If in an iteration there is no packages with empty dependendencies -- report cycle dependency
	result := make([]string, 0, len(packages))
	for len(packageToBeforeIDs) > 0 {
		var ready []string
		for id, deps := range packageToBeforeIDs {
			if len(deps) == 0 {
				ready = append(ready, id)
			}
		}

		if len(ready) == 0 {
			return make([]string, 0), &OrderError{
				what: fmt.Sprintf("Cycle dependency detected after processing %v packages", result),
			}
		}

		for _, id := range ready {
			result = append(result, id)
			delete(packageToBeforeIDs, id)
		}

		for k, v := range packageToBeforeIDs {
			for _, id := range ready {
				index := slices.Index(v, id)
				if index != -1 {
					v = slices.Delete(v, index, index+1)
				}
			}
			packageToBeforeIDs[k] = v
		}
	}

	return result, nil
}

// TODO: how to search, who provides the feature?
// maybe we do not need this for now
