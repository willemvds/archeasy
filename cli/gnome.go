package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"vds.io/archeasy/ansiseq"
	"vds.io/archeasy/jobs"
)

func GnomeSettings(logger *slog.Logger, stdout io.Writer, stderr io.Writer) error {
	ansiseq.TFS_Status(stdout)
	fmt.Fprintf(stdout, "[ * ] Updating GNOME Settings (keybinds).\n")
	ansiseq.Reset(stdout)

	ctx := context.Background()
	err1 := jobs.GnomeKeybind(ctx, logger, "switch-to-workspace-1", "<Super>1")
	if err1 != nil {
		fmt.Fprintf(stdout, "[ F ] (switch-to-workspace-1) %s.\n", err1)
	}

	err2 := jobs.GnomeKeybind(ctx, logger, "switch-to-workspace-2", "<Super>2")
	if err2 != nil {
		fmt.Fprintf(stdout, "[ F ] (switch-to-workspace-2) %s.\n", err2)
	}

	err3 := jobs.GnomeKeybind(ctx, logger, "switch-to-workspace-3", "<Super>3")
	if err3 != nil {
		fmt.Fprintf(stdout, "[ F ] (switch-to-workspace-3) %s.\n", err3)
	}

	err4 := jobs.GnomeKeybind(ctx, logger, "switch-to-workspace-4", "<Super>4")
	if err4 != nil {
		fmt.Fprintf(stdout, "[ F ] (switch-to-workspace-4) %s.\n", err4)
	}

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return errors.New("1 or more keybinds failed")
	}

	ansiseq.RGB(stdout, 20, 210, 10)
	fmt.Fprintf(stdout, "[ OK ] Update GNOME Settings (keybinds).\n")
	ansiseq.Reset(stdout)

	return nil
}
