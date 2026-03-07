package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/toliak/mce/osinfo"
)

func mainInternal() error {
	args, err := ParseArgs(os.Args[1:])
	if err != nil {
		return err
	}

	if !args.CheckDisable {
		if ok, errors := osinfo.CheckPlatform(); !ok {
			fmt.Println("Platform checks failed!")
			for _, error := range errors {
				fmt.Printf("- %s\n", error);
			}
			return fmt.Errorf("Platform checks failed")
		}
	}

	if args.CheckOnly {
		return nil
	}

	info := osinfo.Harvest()

	if args.HarvestOnly {
		info_json, err := json.Marshal(info)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", info_json)
		return nil
	} else {
		fmt.Println("Harvested platform information")
	}

	// err = platform.UpdateRepositories(&info.PkgManager)
	// if err != nil {
	// 	return fmt.Errorf("Update error: %w", err)
	// }
	// found, notFound, err := platform.SearchPackageFullNames(&info.PkgManager, []string{"python3", "zsh"})
	// fmt.Printf("%#v, %#v, %s\n", found, notFound, err)

	// err = platform.InstallPackages(&info.PkgManager, found)
	// if err != nil {
	// 	panic(err)
	// }

	// TODO: use Verbosity
	// args.Verbosity

	fmt.Printf("%s\n", "Finish")
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
