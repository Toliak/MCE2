package main

import (
	"encoding/json"
	"fmt"
	"os"

	// "flag"

	"github.com/toliak/mce/osinfo"
	"github.com/toliak/mce/platform"
)

func main() {
	args, err := ParseArgs(os.Args[1:])
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
		return
	}

	if !args.CheckDisable {
		if ok, errors := osinfo.CheckPlatform(); !ok {
			fmt.Println("Platform checks failed!")
			for _, error := range errors {
				fmt.Printf("- %s\n", error);
				os.Exit(1)
				return
			}
		}
	}

	if args.CheckOnly {
		return
	}

	info := osinfo.Harvest()
	// fmt.Println("Hello, World!")
	// fmt.Printf("%#v\n", info)
	info_json, err := json.Marshal(info)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
		return
	}

	if args.HarvestOnly {
		fmt.Printf("%s\n", info_json)
		return
	} else {
		fmt.Println("Stage2")
	}

	err = platform.UpdateRepositories(&info.PkgManager)
	if err != nil {
		fmt.Printf("update error: %s\n", err)
	}
	found, notFound, err := platform.SearchPackageFullNames(&info.PkgManager, []string{"python3", "zsh"})
	fmt.Printf("%#v, %#v, %s\n", found, notFound, err)

	err = platform.InstallPackages(&info.PkgManager, found)
	if err != nil {
		panic(err)
	}

	// TODO: use Verbosity
	// args.Verbosity
}
