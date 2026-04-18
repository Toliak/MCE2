package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"

	"github.com/toliak/mce/cmd/mce/confirmui"
	"github.com/toliak/mce/cmd/mce/ui"
	"github.com/toliak/mce/inspector"
	"github.com/toliak/mce/platform"
	tb "github.com/toliak/mce/tegnbuilder"
	"github.com/toliak/mce/tegns"
)

func applyPresetToApp(initResult tb.TegnsettInitializeResult, osInfo tb.OSInfoExt, app ui.App, preset JSONPreset) error {
	errorList := make([]string, 0)
	for k, v := range preset {
		if _, ok := initResult.AllIDsSet[k]; !ok {
			errorList = append(errorList, fmt.Sprintf("Tegn or Tegnsett with ID '%s' not found", k))
			continue
		}

		app.State.EnabledIDsMap[k] = v.Enabled
		if v.Params == nil {
			continue
		}
		tegn, ok := initResult.TegnByID[k]
		if !ok {
			errorList = append(errorList, fmt.Sprintf("Unable to set parameter to non-Tegn '%s'", k))
			continue
		}
		params := tegn.GetParameters(osInfo)
		paramsByID := make(map[string]tb.TegnParameter, len(params))
		for _, v := range params {
			paramsByID[v.GetID()] = v
		}
		for pk, pv := range v.Params {
			if _, ok := paramsByID[pk]; !ok {
				errorList = append(errorList, fmt.Sprintf("Parameter with ID '%s' of Tegn '%s'", pk, k))
				continue
			}

			if _, ok := app.State.ParameterByIDMap[k]; !ok {
				app.State.ParameterByIDMap[k] = tb.TegnParameterMap{
					pk: pv,
				}
			} else {
				app.State.ParameterByIDMap[k][pk] = pv
			}
		}
	}

	if len(errorList) != 0 {
		fmt.Println("Preset apply errors:")
		for _, v := range errorList {
			fmt.Printf("- %s\n", v)
		}
		return fmt.Errorf("Preset apply error")
	}

	return nil
}

type TempStateSave struct {
	EnabledIDsMap        tb.TegnGeneralEnabledIDsMap
	ParameterByIDMap     map[string]tb.TegnParameterMap
}

func NewTempStateSave(state *ui.UIState) *TempStateSave {
	return &TempStateSave{
		EnabledIDsMap: state.EnabledIDsMap,
		ParameterByIDMap: state.ParameterByIDMap,
	}
}

func (s TempStateSave) moveIntoState(state *ui.UIState) {
	state.EnabledIDsMap = s.EnabledIDsMap
	state.ParameterByIDMap = s.ParameterByIDMap
}

var stateFilePath string = filepath.Join(os.TempDir(), "mce2-configapp-state-GAaqjz4n.json")

func loadTempState(stateFilePath string) (*TempStateSave, error) {
	byteText, err := os.ReadFile(stateFilePath)
	if err != nil {
		return nil, fmt.Errorf("Unable to read the file '%s': %w\n", stateFilePath, err)
	}
	var dataToLoad TempStateSave
	err = json.Unmarshal(byteText, &dataToLoad)
	if err != nil {
		return nil, fmt.Errorf("Unable to json.Unmarshal the temp state: %w\n", err)
	}

	return &dataToLoad, nil
}

func saveTempState(stateFilePath string, state *ui.UIState) error {
	tempStateDTO := NewTempStateSave(state)
	dataToSave, err := json.Marshal(tempStateDTO)
	if err != nil {
		return fmt.Errorf("json.Marshal error: %w\n", err)
	}
	err = os.WriteFile(stateFilePath, []byte(dataToSave), 0644)
	if err != nil {
		return fmt.Errorf("os.WriteFile error: %w\n", err)
	}

	return nil
}

func prepareApp(args *Args) (*ui.App, error) {
	data, err := inspector.InspectAndHarvest(args.InspectorConfig)
	if err != nil {
		return nil, fmt.Errorf("InspectAndHarvest error: %w", err)
	}
	if data == nil {
		return nil, fmt.Errorf("No harvest data obtained, internal error")
	}

	harvestData := *data
	fmt.Println("Performed checks and harvested platform information")

	if harvestData.OSInfo == nil {
		return nil, fmt.Errorf("Unable to continue without the OSInfo")
	}

	availablePackages := make([]string, 0)
	if harvestData.AvailableManagerPackages != nil {
		availablePackages = *harvestData.AvailableManagerPackages
	}
	
	// List of strings into the map
	availablePackagesMap := make(tb.AvailablePackagesMap)
	for _, v := range availablePackages {
		availablePackagesMap[v] = true
	}

	builderData := tb.OSInfoExt {
		OSInfo: *harvestData.OSInfo,
		AvailableManagerPackages: availablePackagesMap,
		MainInstallDir: args.MainInstallDir,
		DataDir: args.DataDir,
		MceRepositoryURL: args.MceRepositoryURL,
		MceRepositoryBranch: args.MceRepositoryBranch,
	}

	tegnsetts, err := tb.InitializeAllTegnsetts(
		tegns.Tegnsetts,
	)
	if err != nil {
		return nil, fmt.Errorf("InitializeAllTegnsetts error: %w", err)
	}
	initResult := *tegnsetts

	// Just check that we do not have errors
	_, err = tb.GetTegnsettsOrder(initResult.TegnsettByID)
	if err != nil {
		return nil, fmt.Errorf("GetTegnsettsOrder error: %w", err)
	}

	// Installed cache
	alreadyInstalled := make(tb.AvailablePackagesMap)
	alreadyInstalledFeatures := make(tb.TegnInstalledFeaturesMap, len(alreadyInstalled))
	for k, v := range tegnsetts.TegnByID {
		if v.IsInstalled(builderData) {
			alreadyInstalled[k] = true
			features := v.GetFeatures()
			for ft, _ := range features {
				alreadyInstalledFeatures[ft] = true
			}
		}
	}

	app := ui.NewApp(initResult, builderData, alreadyInstalled, alreadyInstalledFeatures)

	for k, tegn := range tegnsetts.TegnByID {
		newParameterMap := make(tb.TegnParameterMap)
		params := tegn.GetParameters(builderData)
		for _, param := range params {
			newParameterMap[param.GetID()] = param.GetDefaultValue()
		}

		app.State.ParameterByIDMap[k] = newParameterMap

	}

	if platform.FileEntryExists(stateFilePath) {
		tempState, err := loadTempState(stateFilePath)
		if err == nil {
			tempState.moveIntoState(app.State)
		} else {
			fmt.Printf("Temp state will not be recovered: %s\n", err)
		}
	}

	err = applyPresetToApp(initResult, builderData, app, args.JSONPreset)
	if err != nil {
		return nil, fmt.Errorf("applyPresetToApp error: %w", err)
	}

	if args.SelectEverything {
		for k := range initResult.TegnByID {
			app.State.EnabledIDsMap[k] = true
		}
		for k := range initResult.TegnsettByID {
			app.State.EnabledIDsMap[k] = true
		}
	}

	return &app, nil
}

func mainInternal() error {
	args, err := ParseArgs(os.Args[1:])
	if err != nil {
		return fmt.Errorf("ParseArgs error: %w", err)
	}

	app, err := prepareApp(args)
	if err != nil {
		return fmt.Errorf("prepareApp error: %w", err)
	}

	if !args.NoUI {
		err = app.Run()
		if err != nil {
			return err
		}
	} else {
		app.State.ExitConfirmed = true
	}

	if !app.State.ExitConfirmed {
		err := saveTempState(stateFilePath, app.State)
		if err != nil {
			fmt.Printf("Unable to save the temp state: %s\n", err)
		}

		return nil
	}

	// TODO: if any during the installation errors -- prompt before continue

	fmt.Printf("App: %#v\n", app)

	// fmt.Printf("%#v\n", tegnsetts)

	// err = RunTUI(tegnsetts, builderData)
	// if err != nil {
	// 	return err
	// }

	tegnsettsObjs := make([]map[string]any, 0, len(app.State.InitResult.TegnsettByID))
	for id, tegnsett := range app.State.InitResult.TegnsettByID {
		children := tegnsett.GetChildren()
		childrenObjs := make([]map[string]any, len(children))
		for j, v := range children {
			params := v.GetParameters(app.State.OSInfExt)
			paramsObjs := make([]map[string]any, len(params))
			for i, p := range params {
				id := v.GetID()
				paramID := p.GetID()
				// paramName := p.GetName()
				paramValue, ok := app.State.ParameterByIDMap[id][paramID]
				if !ok {
					continue
				}
				paramsObjs[i] = map[string]any {
					"id": id,
					"name": p.GetName(),
					"value": paramValue,
					"type": p.GetParamType().String(),
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
	fmt.Println("-----------------------------")
	marshalled, _ = json.MarshalIndent(app.State.EnabledIDsMap, "", "  ")
	fmt.Println(string(marshalled))

	order, err := tb.GetTegnsettsOrder(app.State.InitResult.TegnsettByID)
	if err != nil {
		return err
	}

	availability := tb.GetTegnsettsAvailability(
		app.State.OSInfExt,
		*order,
		app.State.InitResult.TegnsettByID,
		app.State.InitResult.TegnByID,
		app.State.EnabledIDsMap,
		app.State.InstalledFeatures,
	)

	// And list of that into the "to install"
	toInstall := make([]confirmui.ToInstallData, 0)
	for _, id := range order.Tegnsett {
		tegnList := order.TegnByTegnsettID[id]
		
		resultTegnIDList := make([]string, 0, len(tegnList))

		for _, tegnID := range tegnList {
			// tegn := initResult.TegnByID[tegnID]
			// fmt.Printf("%s:%s\n", tegnID, app.State.EnabledIDsMap[tegnID])
			if !app.State.EnabledIDsMap[tegnID] {
				// It is not selected
				continue
			}
			if app.State.InstalledCache[tegnID] {
				fmt.Printf("Selected already installed Tegn '%s', skipped\n", tegnID)
				continue
			}
			if !availability[tegnID].Available {
				fmt.Printf("Selected unavailable Tegn '%s', skipped\n", tegnID)
				continue
			}
			
			resultTegnIDList = append(resultTegnIDList, tegnID)
		}

		toInstall = append(toInstall, confirmui.ToInstallData{
			TegnsettID: id,
			TegnIDList: resultTegnIDList,
		})
	}

	// TODO: split install and the update
	fmt.Println("----------\nWill be installed:")
	for _, d := range toInstall {
		for _, id := range d.TegnIDList {
			fmt.Printf("- %s\n", id)
		}
	}

	// TODO: print also parameters
	installedFeatures := make(tb.TegnInstalledFeaturesMap, len(app.State.InstalledFeatures))
	maps.Copy(installedFeatures, app.State.InstalledFeatures)
	for _, d := range toInstall {
		installedTegns := make([]tb.Tegn, 0, len(d.TegnIDList))

		for _, id := range d.TegnIDList {
			tegn := app.State.InitResult.TegnByID[id]
			err := tegn.ExecInstall(
				app.State.OSInfExt, 
				installedFeatures,
				app.State.ParameterByIDMap[id],
			)
			if err != nil {
				return fmt.Errorf("ExecInstall '%s' error: %w", id, err)
			}
			
			installedTegns = append(installedTegns, tegn)
			maps.Copy(installedFeatures, tegn.GetFeatures())
		}

		err := app.State.InitResult.TegnsettByID[d.TegnsettID].ExecPostInstall(
			installedTegns,
			app.State.OSInfExt,
			installedFeatures,
			app.State.ParameterByIDMap,
		)
		if err != nil {
			return fmt.Errorf("ExecPostInstall '%s' error: %w", d.TegnsettID, err)
		}
	}

	// TODO: add one more app to confirm the
	// TODO: add flag to skip confirmation


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
