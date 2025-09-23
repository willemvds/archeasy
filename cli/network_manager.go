package cli

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"os/user"
)

func InstallNetworkManager(stdout io.Writer, stderr io.Writer) error {
	currentUser, err := user.Current()
	if err != nil || currentUser.Uid != RootId {
		return ErrRootRequired
	}

	fmt.Fprintf(stdout, "[ * ] Installing Network Manager.\n")

	ctx := context.Background()
	installCmd := exec.CommandContext(ctx, "pacman", "-Syy", "--noconfirm", "networkmanager")
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

	fmt.Fprintf(stdout, "[ OK ] Installed NetworkManager package.\n")
	return err
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
