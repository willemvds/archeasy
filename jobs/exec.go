package jobs

import (
	"bytes"
	"context"
	"io"
	osexec "os/exec"
)

// ProgressCallback will always receive the full buffer.
// The buffer should be treated as read-only.
type ProgressCallback func([]byte)

func collect(rdr io.Reader, bufSize uint, doneCh chan<- []byte, progressCallbacks ...ProgressCallback) {
	var all bytes.Buffer
	buf := make([]byte, bufSize)
	for {
		n, err := rdr.Read(buf)
		if n > 0 {
			all.Write(buf[:n])
			for _, cb := range progressCallbacks {
				cb(all.Bytes())
			}
		}
		if err != nil {
			doneCh <- all.Bytes()
			close(doneCh)
			return
		}
	}
}

type ExecJob struct {
	stdoutCallbacks []ProgressCallback
	stderrCallbacks []ProgressCallback
}

func NewExecJob() ExecJob {
	xj := ExecJob{
		stdoutCallbacks: make([]ProgressCallback, 0),
		stderrCallbacks: make([]ProgressCallback, 0),
	}
	return xj
}

func (xj *ExecJob) StdoutCallback(cb ProgressCallback) {
	xj.stdoutCallbacks = append(xj.stdoutCallbacks, cb)
}

func (xj *ExecJob) StderrCallback(cb ProgressCallback) {
	xj.stderrCallbacks = append(xj.stderrCallbacks, cb)
}

func (ej *ExecJob) exec(ctx context.Context, name string, args ...string) (stdout []byte, stderr []byte, err error) {
	cmd := osexec.CommandContext(ctx, name, args...)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, nil, err
	}

	stdoutCh := make(chan []byte, 2)
	stderrCh := make(chan []byte, 2)

	go collect(stdoutPipe, 256, stdoutCh, ej.stdoutCallbacks...)
	go collect(stderrPipe, 256, stderrCh, ej.stderrCallbacks...)

	err = cmd.Wait()
	stdout = <-stdoutCh
	stderr = <-stderrCh
	return stdout, stderr, err
}
