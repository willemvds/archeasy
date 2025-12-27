package jobs

import (
	"context"
	"fmt"
	"log/slog"
)

func GnomeKeybind(ctx context.Context, logger *slog.Logger, action string, key string) error {
	job := NewExecJob()
	job.StdoutCallback(StdoutLogProgressCallbackFn(logger))
	job.StderrCallback(StderrLogProgressCallbackFn(logger))

	_, _, err := job.exec(ctx, "gsettings", "set", "org.gnome.desktop.wm.keybindings", action, fmt.Sprintf("['%s']", key))

	return err
}
