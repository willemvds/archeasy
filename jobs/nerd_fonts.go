package jobs

import (
	"context"
	"log/slog"
)

func InstallNerdFonts(ctx context.Context, logger *slog.Logger) error {
	job := NewExecJob()
	job.StdoutCallback(StdoutLogProgressCallbackFn(logger))
	job.StderrCallback(StderrLogProgressCallbackFn(logger))

	_, _, err := job.exec(ctx, "pacman", "-Syy", "--noconfirm", "nerd-fonts")

	return err
}
