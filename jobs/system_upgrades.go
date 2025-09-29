package jobs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	osexec "os/exec"
	"slices"
	"time"
)

type execResult struct {
	err    error
	stdout []byte
	stderr []byte
}

var ErrPacmanDBLocked = errors.New("pacman DB is locked")

func InstallSystemUpgrades(ctx context.Context, logger *slog.Logger) error {
	job := NewExecJob()
	job.StdoutCallback(StdoutLogProgressCallbackFn(logger))
	job.StderrCallback(StderrLogProgressCallbackFn(logger))
	pidCh := make(chan int, 2)
	job.PidChannel(pidCh)
	allHopeIsLostCh := make(chan error)
	execCh := make(chan execResult)

	go func() {
		pid := <-pidCh
		statPath := fmt.Sprintf("/proc/%d/stat", pid)
		prevStat := []byte("")
		for {
			stat, err := os.ReadFile(statPath)
			if err != nil {
				allHopeIsLostCh <- err
				return
			}
			if slices.Equal(prevStat, stat) {
				fmt.Println(string(prevStat), string(stat))
				allHopeIsLostCh <- errors.New("process might not be making progress...")
				return
			}
			prevStat = stat
			time.Sleep(10 * time.Second)
		}
	}()

	go func() {
		stdout, stderr, err := job.exec(ctx, "pacman", "-Syyu", "--noconfirm")
		execCh <- execResult{
			err:    err,
			stdout: stdout,
			stderr: stderr,
		}
	}()

	select {
	case res := <-execCh:
		if _, ok := res.err.(*osexec.ExitError); ok {
			if bytes.Contains(res.stderr, []byte("unable to lock database")) {
				return ErrPacmanDBLocked
			}
		}
		return res.err
	case err := <-allHopeIsLostCh:
		return err
	}
}
