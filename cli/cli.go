package cli

import (
	"errors"
	"io"
	"log/slog"
)

const RootId = "0"

var ErrRootRequired = errors.New("root required")

func PostInstall(logger *slog.Logger, args []string, stdout io.Writer, stderr io.Writer) error {
	var err error

	err = InstallNetworkManager(logger, stdout, stderr)
	if err != nil {
		return err
	}

	err = StartNetworkManager(stdout, stderr)
	if err != nil {
		return err
	}

	err = installMicrocodePackage(stdout, stderr)
	if err != nil {
		return err
	}

	err = installGraphicsCardDriver(stdout, stderr)
	if err != nil {
		return err
	}

	err = installDesktopEnvironment(stdout, stderr)
	if err != nil {
		return err
	}

	return nil
}
