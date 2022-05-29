// Copyright (C) 2022 The Avly Trader Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package interfaces

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type CmdRunner interface {
	RunCmdSync(cmdLine string, env *[]string) (output string, proc *exec.Cmd, err error)
	RunCmdAsync(cmdLine string, env *[]string) (output string, proc *exec.Cmd, err error)
	PanicCmdSync(cmdLine string, env *[]string) (output string, proc *exec.Cmd)
	PanicCmdAsync(cmdLine string, env *[]string) (output string, proc *exec.Cmd)
}

type SafeCmdRunner struct{}

type SpySafeCmdRunner struct {
	Calls       int
	LastCommand string
}

func (s *SpySafeCmdRunner) RunCmdSync(cmdLine string, env *[]string) (output string, proc *exec.Cmd, err error) {
	s.Calls++
	s.LastCommand = cmdLine

	return cmdLine, nil, nil
}

func (s *SafeCmdRunner) RunCmdSync(cmdLine string, env *[]string) (output string, proc *exec.Cmd, err error) {
	cmdArr := append([]string{"sh", "-c"}, cmdLine)
	var outBuf bytes.Buffer
	output = cmdLine

	executable, err := exec.LookPath(cmdArr[0])
	if err != nil {
		return
	}

	proc = &exec.Cmd{
		Path: executable,
		Args: cmdArr,
		Env:  *env,
		SysProcAttr: &syscall.SysProcAttr{
			Setpgid: true,
		},
		Stdout: &outBuf,
	}

	if errRun := proc.Run(); errRun != nil {
		err = fmt.Errorf("running command \"%s\" not successful: %s", cmdLine, errRun.Error())
	}
	proc.Stdout = nil
	output = strings.TrimSpace(outBuf.String())
	outBuf.Reset()

	return
}

func (s *SpySafeCmdRunner) RunCmdAsync(cmdLine string, env *[]string) (output string, proc *exec.Cmd, err error) {
	return s.RunCmdSync(cmdLine, env)
}

func (s *SafeCmdRunner) RunCmdAsync(cmdLine string, env *[]string) (output string, proc *exec.Cmd, err error) {
	cmdArr := append([]string{"sh", "-c"}, cmdLine)
	var outBuf bytes.Buffer
	output = cmdLine

	executable, err := exec.LookPath(cmdArr[0])
	if err != nil {
		return
	}

	proc = &exec.Cmd{
		Path: executable,
		Args: cmdArr,
		Env:  *env,
		SysProcAttr: &syscall.SysProcAttr{
			Setpgid: true,
		},
		Stdout: &outBuf,
	}

	if errStart := proc.Start(); errStart != nil {
		err = fmt.Errorf("starting command \"%s\" not successful", cmdLine)
		return
	}
	time.Sleep(2 * time.Second)
	output = strings.TrimSpace(outBuf.String())

	return
}

func (s *SpySafeCmdRunner) PanicCmdSync(cmdLine string, env *[]string) (output string, proc *exec.Cmd) {
	s.Calls++
	s.LastCommand = cmdLine

	return cmdLine, nil
}

func (s *SafeCmdRunner) PanicCmdSync(cmdLine string, env *[]string) (output string, proc *exec.Cmd) {
	var err error
	output, proc, err = s.RunCmdSync(cmdLine, env)
	if err != nil {
		panic(err)
	}

	return
}

func (s *SpySafeCmdRunner) PanicCmdAsync(cmdLine string, env *[]string) (output string, proc *exec.Cmd) {
	return s.PanicCmdSync(cmdLine, env)
}

func (s *SafeCmdRunner) PanicCmdAsync(cmdLine string, env *[]string) (output string, proc *exec.Cmd) {
	var err error
	output, proc, err = s.RunCmdAsync(cmdLine, env)
	if err != nil {
		panic(err)
	}

	return
}
