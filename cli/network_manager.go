package cli

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"os/user"

	"vds.io/archeasy/jobs"
)

func InstallNetworkManager(logger *slog.Logger, stdout io.Writer, stderr io.Writer) error {
	currentUser, err := user.Current()
	if err != nil || currentUser.Uid != RootId {
		return ErrRootRequired
	}

	fmt.Fprintf(stdout, "[ * ] Installing Network Manager.\n")

	ctx := context.Background()
	stdoutReceived := 0
	stdoutPrinted := 0
	stdoutHandler := func(p []byte) {
		diff := len(p) - stdoutReceived
		stdoutReceived += diff
		bio := bufio.NewReader(bytes.NewReader(p[stdoutPrinted:]))
		for {
			line, lineErr := bio.ReadBytes('\n')
			if lineErr != nil {
				break
			}
			logger.Info(
				"networkmanager",
				slog.String("stream", "stdout"),
				slog.String("line", string(line)),
			)
			stdoutPrinted += len(line)
		}
	}
	stderrReceived := 0
	stderrPrinted := 0
	stderrHandler := func(p []byte) {
		diff := len(p) - stderrReceived
		stderrReceived += diff
		bio := bufio.NewReader(bytes.NewReader(p[stderrPrinted:]))
		for {
			line, lineErr := bio.ReadBytes('\n')
			if lineErr != nil {
				break
			}
			logger.Info(
				"networkmanager",
				slog.String("stream", "stderr"),
				slog.String("line", string(line)),
			)
			stderrPrinted += len(line)
		}
	}

	err = jobs.InstallNetworkManager(ctx, stdoutHandler, stderrHandler)
	if err != nil {
		fmt.Fprintf(stdout, "[ F ] %s.\n", err)
		return err
	}

	fmt.Fprintf(stdout, "[ OK ] Installed NetworkManager package.\n")
	return nil
}

func StartNetworkManager(stdout io.Writer, stderr io.Writer) error {
	currentUser, err := user.Current()
	if err != nil || currentUser.Uid != RootId {
		return ErrRootRequired
	}

	fmt.Fprintf(stdout, "[ * ] Starting Network Manager.\n")

	ctx := context.Background()
	installCmd := exec.CommandContext(ctx, "systemctl", "start", "NetworkManager")
	stdoutPipe, err := installCmd.StdoutPipe()
	if err != nil {
		return nil
	}
	stderrPipe, err := installCmd.StderrPipe()
	if err != nil {
		return nil
	}

	err = installCmd.Start()
	if err != nil {
		return nil
	}

	go func(stdout io.Reader, cb func([]byte)) {
		buf := make([]byte, 256)
		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				cb(buf[:n])
			}
			if err != nil {
				return
			}
		}
	}(stdoutPipe, func(b []byte) {
		fmt.Fprintf(stdout, "stdout:\t%s", string(b))
	})

	go func(stderr io.Reader, cb func([]byte)) {
		buf := make([]byte, 256)
		for {
			n, err := stderr.Read(buf)
			if n > 0 {
				cb(buf[:n])
			}
			if err != nil {
				return
			}
		}
	}(stderrPipe, func(b []byte) {
		fmt.Fprintf(stderr, "stderr:\t%s", string(b))
	})

	err = installCmd.Wait()

	if err != nil {
		fmt.Fprintf(stdout, "[ F ] %s.\n", err)
		return err
	}

	fmt.Fprintf(stdout, "[ OK ] Started NetworkManager.\n")
	return err
}
