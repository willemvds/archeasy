package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"vds.io/archeasy/cli"
	"vds.io/archeasy/exitcode"
)

func main() {
	logPath := fmt.Sprintf("archeasy-%s.log", time.Now().Format("2006_01_02_1504_05"))
	logfh, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(exitcode.Failure)
	}
	logger := slog.New(slog.NewJSONHandler(logfh, nil))

	args := os.Args
	if len(args) < 2 {
		args = append(args, "help")
	}
	switch args[1] {
	case "post-install":
		fmt.Println(cli.PostInstall(logger, args, os.Stdout, os.Stderr))
	case "networkmanager":
		err := cli.InstallNetworkManager(logger, os.Stdout, os.Stderr)
		if err != nil {
			fmt.Println(err)
			os.Exit(exitcode.Failure)
		}
		fmt.Println(cli.StartNetworkManager(os.Stdout, os.Stderr))
	case "nerdfonts":
		err := cli.InstallNerdFonts(logger, os.Stdout, os.Stderr)
		if err != nil {
			fmt.Println(err)
			os.Exit(exitcode.Failure)
		}
	case "system-upgrades":
		fmt.Println(cli.InstallSystemUpgrades(logger, os.Stdout, os.Stderr))
	default:
		fmt.Println("Usage: archeasy <command>")
		fmt.Println("Commands:")
		fmt.Println("\tpost-install")
		fmt.Println("\tnetworkmanager")
		fmt.Println("\tnerdfonts")
		fmt.Println("\tsystem-upgrades")
		os.Exit(exitcode.Usage)
	}
}
