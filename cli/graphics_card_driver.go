package cli

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"os/user"
)

func installGraphicsCardDriver(stdout io.Writer, stderr io.Writer) error {
	currentUser, err := user.Current()
	if err != nil || currentUser.Uid != RootId {
		return ErrRootRequired
	}

	fmt.Fprintf(stdout, "[ * ] Installing NVIDIA Graphics Card Driver.\n")

	// Error! Bad return status for module build on kernel: 6.16.4-arch1-1 (x86_64)
	// Consult /var/lib/dkms/nvidia/580.76.05/build/make.log for more information.
	// ==> WARNING: `dkms install --no-depmod nvidia/580.76.05 -k 6.16.4-arch1-1' exited 10

	ctx := context.Background()
	installCmd := exec.CommandContext(ctx, "pacman", "-Syy", "--noconfirm", "linux-headers", "linux-lts-headers", "nvidia-open-dkms")
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

	fmt.Fprintf(stdout, "[ OK ] Installed NVIDIA Graphics Card Driver package.\n")
	return err
}
