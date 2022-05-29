// Copyright (C) 2022 The Avly Trader Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package helpers

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	ifc "github.com/9tmark/avly-trader/internal/interfaces"
)

type ProcDeathQueue []*exec.Cmd

func RestOrDie(runner ifc.CmdRunner, env *[]string, procs ...*exec.Cmd) {
	for i := 0; i < len(procs); i++ {
		var logCommand string
		restStatus := make(chan bool)

		go hasFoundRest(procs[i], restStatus)

		select {
		case hasRest := <-restStatus:
			if hasRest {
				logCommand = fmt.Sprintf("echo $(date +\"%%Y/%%m/%%d %%T\") <%d>'s soul was calmed. It left free memory: %s >> $AVL_LOGS/zombie.log", procs[i].ProcessState.Pid(), procs[i].Path)
				runner.RunCmdSync(logCommand, env)
				continue
			}
			// syscall.Kill(-procs[i].SysProcAttr.Pgid, syscall.SIGKILL)
			procs[i].Process.Kill()
			logCommand = fmt.Sprintf("echo $(date +\"%%Y/%%m/%%d %%T\") <%d> wasn't afraid about death, they just didn't want to be there when it happened: %s >> $AVL_LOGS/zombie.log", procs[i].ProcessState.Pid(), procs[i].Path)
			runner.RunCmdSync(logCommand, env)
		case <-time.After(45 * time.Second):
			// syscall.Kill(-procs[i].SysProcAttr.Pgid, syscall.SIGKILL)
			procs[i].Process.Kill()
			logCommand = fmt.Sprintf("echo $(date +\"%%Y/%%m/%%d %%T\") <%d> was so lonely: %s >> $AVL_LOGS/zombie.log", procs[i].ProcessState.Pid(), procs[i].Path)
			runner.RunCmdSync(logCommand, env)
			continue
		}
	}
}

func hasFoundRest(proc *exec.Cmd, evalCh chan bool) {
	if waitErr := proc.Wait(); waitErr != nil {
		if !strings.Contains(waitErr.Error(), "already called") {
			evalCh <- false
			return
		}

		evalCh <- true
		return
	}

	evalCh <- true
	return
}

func (pdq ProcDeathQueue) Add(proc *exec.Cmd) {
	pdq = append(pdq, proc)
}

func (pdq ProcDeathQueue) clear() {
	pdq = nil
	pdq = ProcDeathQueue{}
}

func (pdq ProcDeathQueue) LetDie(runner ifc.CmdRunner, env *[]string) {
	defer pdq.clear()
	RestOrDie(runner, env, pdq...)
}
