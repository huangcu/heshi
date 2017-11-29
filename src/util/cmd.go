package util

// import (
// 	"bufio"
// 	"bytes"
// 	"errors"
// 	"fmt"
// 	"io"
// 	"os"
// 	"os/exec"
// 	"strings"
// 	"time"
// )

// func Run(c string, params ...string) (string, error) {
// 	TraceCmd(c, params...)
// 	out, err := exec.Command(c, params...).CombinedOutput()
// 	msg := fmt.Sprintf("error to run %s %v", c, params)
// 	o := strings.TrimSpace(string(out))
// 	if o != "" {
// 		msg += "; " + o
// 	}
// 	if err != nil {
// 		return "", errors.Wrap(err, msg)
// 	}
// 	return o, nil
// }

// func traceRC(r io.ReadCloser, ec chan error) {
// 	in := bufio.NewScanner(r)
// 	for in.Scan() {
// 		Traceln(in.Text())
// 	}
// 	if err := in.Err(); err != nil {
// 		ec <- err
// 	}
// }

// func RunWithStdOutput(c string, params ...string) error {
// 	TraceCmd(c, params...)
// 	cmd := exec.Command(c, params...)
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	if err := cmd.Run(); err != nil {
// 		return errors.Wrapf(err, "error to run cmd [%s]", getCmd(c, params...))
// 	}
// 	return nil
// }

// func RunWithTrace(c string, params ...string) error {
// 	TraceCmd(c, params...)
// 	cmd := exec.Command(c, params...)
// 	ec := make(chan error)

// 	op, err := cmd.StdoutPipe()
// 	if err != nil {
// 		return err
// 	}
// 	go traceRC(op, ec)

// 	ep, err := cmd.StderrPipe()
// 	if err != nil {
// 		return err
// 	}
// 	go traceRC(ep, ec)

// 	go func() { ec <- cmd.Run() }()

// 	if err := <-ec; err != nil {
// 		return errors.Wrapf(err, "error to run cmd[%s]", getCmd(c, params...))
// 	}
// 	return nil
// }

// func RunCmdInTime(t int, cmd *exec.Cmd) (string, string, error) {
// 	var stdout, stderr bytes.Buffer
// 	cmd.Stdout, cmd.Stderr = &stdout, &stderr

// 	if err := cmd.Start(); err != nil {
// 		return "", "", err
// 	}

// 	ec := make(chan error)
// 	go func() {
// 		ec <- cmd.Wait()
// 	}()

// 	select {
// 	case err := <-ec:
// 		o := strings.TrimSpace(stdout.String())
// 		e := strings.TrimSpace(stderr.String())
// 		return o, e, err
// 	case <-time.After(time.Duration(t) * time.Second):
// 		cmd.Process.Kill()
// 		err := fmt.Errorf("timeout to run cmd [%s] in %ds", cmd.Path, t)
// 		return "", "", err
// 	}
// }

// func RunCmd(cmd *exec.Cmd) (string, string, error) {
// 	var stdout, stderr bytes.Buffer
// 	cmd.Stdout, cmd.Stderr = &stdout, &stderr

// 	err := cmd.Run()

// 	o := strings.TrimSpace(stdout.String())
// 	e := strings.TrimSpace(stderr.String())
// 	return o, e, err
// }

// func RunCmdWithStdout(cmd *exec.Cmd) error {
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	return cmd.Run()
// }
