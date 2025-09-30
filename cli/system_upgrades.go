package cli

import (
	"context"
	"fmt"
	"log/slog"
	"os/user"

	"vds.io/archeasy"
	"vds.io/archeasy/ansiseq"
	"vds.io/archeasy/jobs"
)

func InstallSystemUpgrades(logger *slog.Logger, stdout archeasy.BufferedWriter, stderr archeasy.BufferedWriter) error {
	currentUser, err := user.Current()
	if err != nil || currentUser.Uid != RootId {
		return ErrRootRequired
	}

	ansiseq.TFS_Status(stdout)
	fmt.Fprintf(stdout, "[ * ] Upgrading system packages...")
	ansiseq.Reset(stdout)
	stdout.Flush()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = jobs.InstallSystemUpgrades(ctx, logger)
	if err != nil {
		ansiseq.ClearLine(stdout)
		fmt.Fprintf(stdout, "[ F ] %s.\n", err)
		stdout.Flush()
		return err
	}

	ansiseq.ClearLine(stdout)
	ansiseq.RGB(stdout, 20, 210, 10)
	fmt.Fprintf(stdout, "[ OK ] System packages upgraded.\n")
	ansiseq.Reset(stdout)
	stdout.Flush()

	return nil
}
