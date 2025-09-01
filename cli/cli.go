package cli

import (
	"errors"
	"io"
)

const RootId = "0"

var ErrRootRequired = errors.New("root required")

func PostInstall(args []string, stdout io.Writer, stderr io.Writer) error {
	var err error

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
