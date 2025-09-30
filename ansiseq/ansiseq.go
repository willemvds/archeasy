package ansiseq

import (
	"fmt"
	"io"
)

func Reset(writer io.Writer) {
	fmt.Fprintf(writer, "\033[0m")
}

func TFS_Status(writer io.Writer) {
	fmt.Fprintf(writer, "\033[38;5;81m")
}

func TFS_OK(writer io.Writer) {
	fmt.Fprintf(writer, "\033[38;5;118m")
}

func TFS_Fail(writer io.Writer) {
	fmt.Fprintf(writer, "\033[38;5;196m")
}

func RGB(writer io.Writer, r uint8, g uint8, b uint8) {
	fmt.Fprintf(writer, "\033[38;2;%d;%d;%dm", r, g, b)
}

func ClearLine(writer io.Writer) {
	fmt.Fprintf(writer, "\033[2K\r")
}
