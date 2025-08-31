package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

const RootId = "0"

var ErrRootRequired = errors.New("root required")

type cpuInfo struct {
	vendorId  string
	cpuFamily string
	model     string
	modelName string
	microCode string
}

var availableMicrocodePackages = map[string]string{
	"AuthenticAMD": "amd-ucode",
}

func fetchCpuInfo() (*cpuInfo, error) {
	cpuinfo, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return nil, err
	}

	fetchedCpuInfo := cpuInfo{}

	for line := range strings.SplitSeq(string(cpuinfo), "\n") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "vendor_id":
			fetchedCpuInfo.vendorId = value
		case "cpu family":
			fetchedCpuInfo.cpuFamily = value
		case "model":
			fetchedCpuInfo.model = value
		case "model name":
			fetchedCpuInfo.modelName = value
		case "microcode":
			fetchedCpuInfo.microCode = value
		}
	}

	return &fetchedCpuInfo, nil
}

func installMicrocodePackage(stdout io.Writer, stderr io.Writer) error {
	fmt.Fprintf(stdout, "[   ] Checking CPU µCode package.\n")
	cpuinfo, err := fetchCpuInfo()
	if err != nil {
		return err
	}
	fmt.Println(cpuinfo)

	ucodePkg, exists := availableMicrocodePackages[cpuinfo.vendorId]
	fmt.Println(ucodePkg, exists)
	if !exists {
		return nil
	}
	fmt.Fprintf(stdout, "[   ] Found CPU µCode package for %s.\n", cpuinfo.vendorId)

	currentUser, err := user.Current()
	if err != nil || currentUser.Uid != RootId {
		return ErrRootRequired
	}

	fmt.Fprintf(stdout, "[ * ] Installing CPU µCode package for %s.\n", cpuinfo.vendorId)
	ctx := context.Background()
	installCmd := exec.CommandContext(ctx, "pacman", "-Syy", "--noconfirm", ucodePkg)
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

	fmt.Fprintf(stdout, "[ OK ] Installed CPU µCode package.\n")
	return err
}

func PostInstall(args []string, stdout io.Writer, stderr io.Writer) error {
	err := installMicrocodePackage(stdout, stderr)
	if err != nil {
		return err
	}

	return nil
}
