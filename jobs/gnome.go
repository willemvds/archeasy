package jobs

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

const GnomeClockFormat12H = "12h"
const GnomeClockFormat24H = "24h"

var ErrInvalidClockFormat = errors.New("invalid clock format")

///org/gnome/desktop/interface/clock-format
//  '24h'

// /org/gtk/settings/file-chooser/clock-format
//
//	'24h'
func GnomeClockFormat(ctx context.Context, logger *slog.Logger, value string) error {
	if value != GnomeClockFormat12H && value != GnomeClockFormat24H {
		return ErrInvalidClockFormat
	}

	job := NewExecJob()
	job.StdoutCallback(StdoutLogProgressCallbackFn(logger))
	job.StderrCallback(StderrLogProgressCallbackFn(logger))

	_, _, err := job.exec(ctx, "gsettings", "set", "org.gnome.desktop.interface", "clock-format", value)

	return err
}

func GnomeKeybind(ctx context.Context, logger *slog.Logger, action string, key string) error {
	job := NewExecJob()
	job.StdoutCallback(StdoutLogProgressCallbackFn(logger))
	job.StderrCallback(StderrLogProgressCallbackFn(logger))

	_, _, err := job.exec(ctx, "gsettings", "set", "org.gnome.desktop.wm.keybindings", action, fmt.Sprintf("['%s']", key))

	return err
}
