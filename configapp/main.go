package main

import (
	"fmt"
	"os"
	// "flag"

	"github.com/toliak/mce/osinfo"
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
	fmt.Printf("%#v\n", info)

	if args.HarvestOnly {
		return
	}

	// TODO: use Verbosity
	// args.Verbosity
}
