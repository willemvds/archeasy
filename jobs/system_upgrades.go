package jobs

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"time"
)

func InstallSystemUpgrades(ctx context.Context, logger *slog.Logger) error {
	job := NewExecJob()
	job.StdoutCallback(StdoutLogProgressCallbackFn(logger))
	job.StderrCallback(StderrLogProgressCallbackFn(logger))
	pidCh := make(chan int, 2)
	job.PidChannel(pidCh)
	allHopeIsLostCh := make(chan error)
	execCh := make(chan error)

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
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		_, _, err := job.exec(ctx, "pacman", "-Syyu", "--noconfirm")
		execCh <- err
	}()

	select {
	case err := <-execCh:
		return err
	case err := <-allHopeIsLostCh:
		return err
	}
}
