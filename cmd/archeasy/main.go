package main

import (
	"fmt"
	"os"

	"vds.io/archeasy/cli"
	"vds.io/archeasy/exitcode"
)

func main() {
	args := os.Args
	switch args[1] {
	case "post-install":
		fmt.Println(cli.PostInstall(args, os.Stdout, os.Stderr))
	case "networkmanager":
		err := cli.InstallNetworkManager(os.Stdout, os.Stderr)
		if err != nil {
			os.Exit(exitcode.Failure)
		}
		fmt.Println(cli.StartNetworkManager(os.Stdout, os.Stderr))
	default:
		fmt.Println("Usage: archeasy <command>")
		fmt.Println("Commands:")
		fmt.Println("\tpost-install")
		fmt.Println("\tnetworkmanager")
		os.Exit(exitcode.Usage)
	}
}
