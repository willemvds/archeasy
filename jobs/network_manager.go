package jobs

import (
	"context"
)

func InstallNetworkManager(ctx context.Context, stdoutCb ProgressCallback, stderrCb ProgressCallback) error {
	job := NewExecJob()
	job.StdoutCallback(stdoutCb)
	job.StderrCallback(stderrCb)

	_, _, err := job.exec(ctx, "pacman", "-Syy", "--noconfirm", "networkmanager")

	return err
}
