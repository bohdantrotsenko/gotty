package localcommand

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/creack/pty"
)

const (
	DefaultCloseSignal  = syscall.SIGINT
	DefaultCloseTimeout = 10 * time.Second
)

type LocalCommand struct {
	command string
	argv    []string

	closeSignal  syscall.Signal
	closeTimeout time.Duration

	cmd       *exec.Cmd
	pty       *os.File
	ptyClosed chan struct{}
}

func New(ctx context.Context, command string, argv []string, options ...Option) (*LocalCommand, error) {
	cmd := exec.CommandContext(ctx, command, argv...)

	pty, err := pty.Start(cmd)
	if err != nil {
		// todo close cmd?
		return nil, fmt.Errorf("failed to start command `%s`: %w", command, err)
	}
	ptyClosed := make(chan struct{})

	lcmd := &LocalCommand{
		command: command,
		argv:    argv,

		closeSignal:  DefaultCloseSignal,
		closeTimeout: DefaultCloseTimeout,

		cmd:       cmd,
		pty:       pty,
		ptyClosed: ptyClosed,
	}

	for _, option := range options {
		option(lcmd)
	}

	// When the process is closed by the user,
	// close pty so that Read() on the pty breaks with an EOF.
	go func() {
		defer func() {
			lcmd.pty.Close()
			close(lcmd.ptyClosed)
		}()

		lcmd.cmd.Wait()
	}()

	return lcmd, nil
}

func (lcmd *LocalCommand) Read(p []byte) (n int, err error) {
	return lcmd.pty.Read(p)
}

func (lcmd *LocalCommand) Write(p []byte) (n int, err error) {
	return lcmd.pty.Write(p)
}

func (lcmd *LocalCommand) Close() error {
	if lcmd.cmd != nil && lcmd.cmd.Process != nil {
		lcmd.cmd.Process.Signal(lcmd.closeSignal)
	}
	timeout := lcmd.closeTimeoutC()
	for {
		select {
		case <-lcmd.ptyClosed:
			return nil
		case <-timeout:
			lcmd.cmd.Process.Signal(syscall.SIGKILL)
		}
	}
}

func (lcmd *LocalCommand) WindowTitleVariables() map[string]any {
	return map[string]any{
		"command": lcmd.command,
		"argv":    lcmd.argv,
		"pid":     lcmd.cmd.Process.Pid,
	}
}

func (lcmd *LocalCommand) ResizeTerminal(width int, height int) error {
	return pty.Setsize(lcmd.pty, &pty.Winsize{
		Cols: uint16(width),
		Rows: uint16(height),
	})
}

func (lcmd *LocalCommand) closeTimeoutC() <-chan time.Time {
	if lcmd.closeTimeout >= 0 {
		return time.After(lcmd.closeTimeout)
	}

	return make(chan time.Time)
}
