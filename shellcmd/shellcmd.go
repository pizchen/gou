package shellcmd

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type ShellCmd struct {
	*exec.Cmd
}

var (
	SHELL string = "/bin/sh"
)

func init() {
	if sh, exists := os.LookupEnv("SHELL"); exists {
		SHELL = sh
	}
}

func NewShellCmd(shcmd string, args ...string) *ShellCmd {
	cmdStr := shcmd + " " + strings.Join(args, " ")
	cmd := exec.Command(SHELL, "-c", cmdStr)
	return &ShellCmd{cmd}
}

func (shc *ShellCmd) StdOutErrPipe() (io.ReadCloser, error) {
	if shc.Stdout != nil || shc.Stderr != nil {
		return nil, errors.New("exec: Stdout/Stderr already set")
	}
	if shc.Process != nil {
		return nil, errors.New("exec: StdOutErrPipe after process started")
	}
	pr, pw, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	shc.Stdout = pw
	shc.Stderr = pw
	return pr, nil
}

func (shc *ShellCmd) RunTimeout(sec uint) (err error) {

	err = nil
	if err = shc.Start(); err != nil {
		return err
	}

	done := make(chan error, 1)
	go func() {
		done <- shc.Wait()
	}()

	select {
	case <-time.After(time.Duration(sec) * time.Second):
		err = shc.Process.Kill()
	case err = <-done:
	}
	return
}
