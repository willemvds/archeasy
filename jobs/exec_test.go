package jobs

import (
	"bytes"
	"context"
	osexec "os/exec"
	"reflect"
	"testing"
	"time"
)

func TestExecJob(t *testing.T) {
	job := NewExecJob()
	ctx := context.Background()
	echoText := "Hello, World"

	stdout, stderr, err := job.exec(ctx, "echo", echoText)

	if err != nil {
		t.Errorf("Expected err=nil, got err=%v", err)
	}

	if !bytes.Equal(bytes.TrimRight(stdout, "\n"), []byte("Hello, World")) {
		t.Errorf("Expected stdout='%s', got stdout='%s'", echoText, stdout)
	}

	if len(stderr) != 0 {
		t.Errorf("Expected stderr='', got stderr='%s'", stderr)
	}
}

func TestExecJobTimeout(t *testing.T) {
	job := NewExecJob()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	stdout, stderr, err := job.exec(ctx, "cat", "/dev/urandom")

	if exitErr, isExitErr := err.(*osexec.ExitError); isExitErr {
		if exitErr.Exited() {
			t.Errorf("Expected exited=false, got exited=true")
		}
	} else {
		t.Errorf("Expected err=%t, got err=%v %s", err, err, reflect.TypeOf(err))
	}

	if len(stdout) == 0 {
		t.Errorf("Expected len(stdout) > 0")
	}

	if len(stderr) != 0 {
		t.Errorf("Expected stderr='', got stderr='%s'", stderr)
	}
}

func TestExecJobStdoutCallback(t *testing.T) {
	job := NewExecJob()
	var buf bytes.Buffer
	received := 0
	job.StdoutCallback(func(b []byte) {
		diff := len(b) - received
		buf.Write(b[received:])
		received = received + diff
	})
	ctx := context.Background()
	stdout, stderr, err := job.exec(ctx, "cat", "exec_test.go")

	if err != nil {
		t.Errorf("Expected err=nil, got err=%s", err)
	}

	if len(stderr) != 0 {
		t.Errorf("Expected len(stderr)=0, got len(stderr)=%d", len(stderr))
	}

	if len(stdout) != received || len(stdout) != buf.Len() {
		t.Errorf("Something's gone wrong. len(stdout)=%d, len(callback)=%d, received(callback)=%d", len(stdout), buf.Len(), received)
	}
}
