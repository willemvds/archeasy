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

	keybinds := map[string]string{
		"switch-to-workspace-1": "<Super>1",
		"switch-to-workspace-2": "<Super>2",
		"switch-to-workspace-3": "<Super>3",
		"switch-to-workspace-4": "<Super>4",
		"move-to-workspace-1":   "<Super><Shift>1",
		"move-to-workspace-2":   "<Super><Shift>2",
		"move-to-workspace-3":   "<Super><Shift>3",
		"move-to-workspace-4":   "<Super><Shift>4",
	}

	ctx := context.Background()
	hasFailures := false
	for action, bind := range keybinds {
		err := jobs.GnomeKeybind(ctx, logger, action, bind)
		if err != nil {
			hasFailures = true
			fmt.Fprintf(stdout, "[ F ] (%s->%s) %s.\n", action, bind, err)
		}
	}

	if hasFailures {
		return errors.New("1 or more keybinds failed")
	}

	err := jobs.GnomeClockFormat(ctx, logger, jobs.GnomeClockFormat24H)
	if err != nil {
		fmt.Fprintf(stdout, "[ F ] (Clock Format -> 24h) %s.\n", err)
		return err
	}

	ansiseq.RGB(stdout, 20, 210, 10)
	fmt.Fprintf(stdout, "[ OK ] Update GNOME Settings (keybinds).\n")
	ansiseq.Reset(stdout)

	return nil
}
