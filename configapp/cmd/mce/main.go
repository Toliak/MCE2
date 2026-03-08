package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/toliak/mce/inspector"
	"github.com/toliak/mce/tegnbuilder"
	"github.com/toliak/mce/tegns"
)

func mainInternal() error {
	args, err := ParseArgs(os.Args[1:])
	if err != nil {
		return err
	}

	data, err := inspector.InspectAndHarvest(args.InspectorConfig)
	if err != nil {
		return err
	}

	if data == nil {
		return fmt.Errorf("No harvest data obtained, internal error")
	}

	harvestData := *data

	fmt.Println("Performed checks and harvested platform information")

	if harvestData.OSInfo == nil {
		return fmt.Errorf("Unable to continue without the OSInfo")
	}

	availablePackages := make([]string, 0)
	if harvestData.AvailableManagerPackages != nil {
		availablePackages = *harvestData.AvailableManagerPackages
	}

	builderData := tegnbuilder.TegnBuilderData {
		OSInfo: *harvestData.OSInfo,
		AvailableManagerPackages: availablePackages,
	}

	tegnsetts := tegns.InitializeAllTegnsetts(
		tegns.Tegnsetts,
		builderData,
	)

	// TODO: initialize Tegns

	// Tegn -- package
	// Tegnsett -- category

	// fmt.Printf("%#v\n", tegnsetts)

	tegnsettsObjs := make([]map[string]any, len(tegnsetts))
	for i, tegnsett := range tegnsetts {
		children := tegnsett.GetChildren()
		childrenObjs := make([]map[string]any, len(children))
		for j, v := range tegnsett.GetChildren() {
			params := v.GetParameters()
			paramsObjs := make([]map[string]any, len(params))
			for i, v := range params {
				paramsObjs[i] = map[string]any {
					"name": v.Name,
					"value": v.Value,
					"type": v.ParamType.String(),
				}
			}

			childrenObjs[j] = map[string]any{
				"id":     v.GetID(),
				"name":   v.GetName(),
				"params": paramsObjs,
			}
		}

		tegnsettsObjs[i] = map[string]any{
			"id":       tegnsett.GetID(),
			"name":     tegnsett.GetName(),
			"features": tegnsett.GetFeatures(),
			"children": childrenObjs,
		}
	}

	marshalled, _ := json.Marshal(tegnsettsObjs)
	fmt.Println(string(marshalled))

	return nil
}

func main() {
	err := mainInternal()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
		return
	}
}
