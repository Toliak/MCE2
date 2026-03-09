package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/toliak/mce/cmd/mce/ui"
	"github.com/toliak/mce/inspector"
	tb "github.com/toliak/mce/tegnbuilder"
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

	builderData := tb.TegnBuilderData {
		OSInfo: *harvestData.OSInfo,
		AvailableManagerPackages: availablePackages,
	}

	tegnsetts, err := tb.InitializeAllTegnsetts(
		tegns.Tegnsetts,
		builderData,
	)
	if err != nil {
		return err
	}
	initResult := *tegnsetts

	// Just check that we do not have errors
	_, err = tb.GetTegnsettsOrder(initResult.TegnsettByID)
	if err != nil {
		return err
	}

	app := ui.NewApp(initResult, harvestData)

	for k, v := range DefaultEnables {
		_, ok := app.State.EnabledIDsMap[k]
		if !ok {
			fmt.Println("")
			continue
		}

		app.State.EnabledIDsMap[k] = v
	}

	// TODO: if any errors -- prompt before continue
	
	err = app.Run()
	if err != nil {
		return err
	}

	fmt.Printf("App: %#v\n", app)

	// fmt.Printf("%#v\n", tegnsetts)

	// err = RunTUI(tegnsetts, builderData)
	// if err != nil {
	// 	return err
	// }

	tegnsettsObjs := make([]map[string]any, 0, len(initResult.TegnsettByID))
	for id, tegnsett := range initResult.TegnsettByID {
		children := tegnsett.GetChildren()
		childrenObjs := make([]map[string]any, len(children))
		for j, v := range children {
			params := v.GetParameters()
			paramsObjs := make([]map[string]any, len(params))
			for i, v := range params {
				paramsObjs[i] = map[string]any {
					"name": v.Name,
					"value": v.GetValue(),
					"type": v.ParamType.String(),
				}
			}

			childrenObjs[j] = map[string]any{
				"id":     v.GetID(),
				"name":   v.GetName(),
				"params": paramsObjs,
			}
		}

		tegnsettsObjs = append(tegnsettsObjs, map[string]any{
			"id":       id,
			"name":     tegnsett.GetName(),
			"children": childrenObjs,
		})
	}

	marshalled, _ := json.MarshalIndent(tegnsettsObjs, "", "  ")
	fmt.Println(string(marshalled))
	marshalled, _ = json.MarshalIndent(app.State.EnabledIDsMap, "", "  ")
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
