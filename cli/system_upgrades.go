package cli

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os/user"

	"vds.io/archeasy/ansiseq"
	"vds.io/archeasy/jobs"
)

func InstallSystemUpgrades(logger *slog.Logger, stdout io.Writer, stderr io.Writer) error {
	currentUser, err := user.Current()
	if err != nil || currentUser.Uid != RootId {
		return ErrRootRequired
	}

	ansiseq.TFS_Status(stdout)
	fmt.Fprintf(stdout, "[ * ] Upgrading system packages.\n")
	ansiseq.Reset(stdout)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = jobs.InstallSystemUpgrades(ctx, logger)
	if err != nil {
		fmt.Fprintf(stdout, "[ F ] %s.\n", err)
		return err
	}

	ansiseq.RGB(stdout, 20, 210, 10)
	fmt.Fprintf(stdout, "[ OK ] System packages upgraded.\n")
	ansiseq.Reset(stdout)

	return nil
}
