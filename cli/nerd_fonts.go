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

func InstallNerdFonts(logger *slog.Logger, stdout io.Writer, stderr io.Writer) error {
	currentUser, err := user.Current()
	if err != nil || currentUser.Uid != RootId {
		return ErrRootRequired
	}

	ansiseq.TFS_Status(stdout)
	fmt.Fprintf(stdout, "[ * ] Installing Nerd Fonts.\n")
	ansiseq.Reset(stdout)

	ctx := context.Background()
	err = jobs.InstallNerdFonts(ctx, logger)
	if err != nil {
		fmt.Fprintf(stdout, "[ F ] %s.\n", err)
		return err
	}

	ansiseq.RGB(stdout, 20, 210, 10)
	fmt.Fprintf(stdout, "[ OK ] Installed 'nerd-fonts' package.\n")
	ansiseq.Reset(stdout)

	return nil
}
