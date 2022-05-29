// Copyright (C) 2022 The Avly Trader Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"fmt"
	"time"

	hlp "github.com/9tmark/avly-trader/internal/helpers"
	ifc "github.com/9tmark/avly-trader/internal/interfaces"
)

type FlagInfo struct {
	p      *bool
	fName  string
	sName  string
	defVal bool
	usage  string
}

var env = []string{
	"USER=root",
	"AVL_LOGS=/var/log/avly-trader",
	"THIRD_PARTY=/opt/third-party",
	"WINEPREFIX=/opt/.mtprfx",
	"WINEDEBUG=-all",
	"DISPLAY=:1",
	"SCREEN_NUM=0",
	"SCREEN_WHD=1366x768x16",
	"PATH=/usr/local/bin:/usr/bin:/usr/local/sbin:/usr/sbin:/opt/avly-trader/bin",
}

func main() {
	var isPrepare, isFledge, isLaunch, isStop, isDrain, isCleanUp, isEnter, isMute bool
	mp := &ifc.FmtMsgPrinter{}
	lp := &ifc.LogMsgPrinter{}
	runner := &ifc.SafeCmdRunner{}
	flags := []FlagInfo{
		// verbs
		{p: &isPrepare, fName: "prepare", sName: "p", defVal: false, usage: "verify perquisites for a workstation to work properly"},
		{p: &isFledge, fName: "fledge", sName: "f", defVal: false, usage: "(safely) pull up VNC server"},
		{p: &isLaunch, fName: "launch", sName: "l", defVal: false, usage: "(safely) launch target executable"},
		{p: &isStop, fName: "stop", sName: "s", defVal: false, usage: "stop target process"},
		{p: &isDrain, fName: "drain", sName: "d", defVal: false, usage: "shut down VNC server"},
		{p: &isCleanUp, fName: "clean-up", sName: "c", defVal: false, usage: "dispose remains of target process"},
		{p: &isEnter, fName: "enter", sName: "e", defVal: false, usage: "run startup routine as container process"},
		// options
		{p: &isMute, fName: "mute", sName: "m", defVal: false, usage: "mute output unless error occurs"}, // not supported yet
	}
	verbs := []*bool{&isPrepare, &isFledge, &isLaunch, &isStop, &isDrain, &isCleanUp, &isEnter}
	opts := []*bool{&isMute}

	for i := 0; i < len(flags); i++ {
		v := flags[i]
		flag.BoolVar(v.p, v.fName, v.defVal, v.usage)
		flag.BoolVar(v.p, v.sName, v.defVal, "")
	}

	flag.Parse()
	mp.Printfln("Avly Trader | Cloud Trading CLI")

	switch true {
	default:
	case !hlp.HasOnlyOneTrueValue(verbs...):
		mp.Errorfln("avly: exactly one verb flag is required\nRun with '--help' for usage")
	case isPrepare:
		prepareHandler(mp, lp, runner, opts...)
	case isFledge:
		fledgeHandler(mp, lp, runner, opts...)
	case isLaunch:
		launchHandler(mp, lp, runner, opts...)
	case isCleanUp:
		cleanUpHandler(mp, lp, runner, opts...)
	case isEnter:
		enterHandler(mp, lp, runner, opts...)
	}
}

func prepareHandler(msgPrinter ifc.MsgPrinter, logPrinter ifc.MsgPrinter, runner ifc.CmdRunner, opts ...*bool) {
	if !hlp.WasRunAsRoot(runner) {
		msgPrinter.Errorfln("avly: flag 'prepare' needs to be executed as root")
	}
	_, _, err := prepare(msgPrinter, logPrinter, runner)
	if err != nil {
		logPrinter.Errorfln("avly: %s", err.Error())
	}
}

func fledgeHandler(msgPrinter ifc.MsgPrinter, logPrinter ifc.MsgPrinter, runner ifc.CmdRunner, opts ...*bool) {
	if !hlp.WasRunAsRoot(runner) {
		msgPrinter.Errorfln("avly: flag 'fledge' needs to be executed as root")
	}
	framebufferAlive, vncServerAlive, err := fledge(msgPrinter, logPrinter, runner)
	if err != nil {
		logPrinter.Errorfln("avly: %s", err.Error())
	}
	if !framebufferAlive {
		logPrinter.Errorfln("avly: could not open or verify framebuffer")
	}
	if !vncServerAlive {
		logPrinter.Errorfln("avly: could not pull up or verify VNC server")
	}
}

func launchHandler(msgPrinter ifc.MsgPrinter, logPrinter ifc.MsgPrinter, runner ifc.CmdRunner, opts ...*bool) {
	if !hlp.WasRunAsRoot(runner) {
		msgPrinter.Errorfln("avly: flag 'launch' needs to be executed as root")
	}
	targetProcessAlive, err := launch(msgPrinter, logPrinter, runner)
	if err != nil {
		logPrinter.Errorfln("avly: %s", err.Error())
	}
	if !targetProcessAlive {
		logPrinter.Errorfln("avly: could not launch or verify target executable")
	}
}

func cleanUpHandler(msgPrinter ifc.MsgPrinter, logPrinter ifc.MsgPrinter, runner ifc.CmdRunner, opts ...*bool) {
	if !hlp.WasRunAsRoot(runner) {
		msgPrinter.Errorfln("avly: flag 'clean-up' needs to be executed as root")
	}
	cleanedUp, err := cleanUp(msgPrinter, logPrinter, runner)
	if err != nil {
		logPrinter.Errorfln("avly: %s", err.Error())
	}
	if !cleanedUp {
		logPrinter.Printfln("avly: warn: could not clean up")
	} else {
		logPrinter.Printfln("Cleanup: OK")
	}
}

func stopHandler(msgPrinter ifc.MsgPrinter, logPrinter ifc.MsgPrinter, runner ifc.CmdRunner, opts ...*bool) {
	if !hlp.WasRunAsRoot(runner) {
		msgPrinter.Errorfln("avly: flag 'stop' needs to be executed as root")
	}
	targetProcessDead, err := stop(msgPrinter, logPrinter, runner)
	if err != nil {
		logPrinter.Errorfln("avly: %s", err.Error())
	}
	if !targetProcessDead {
		logPrinter.Errorfln("avly: could not stop target process")
	}
}

func drainHandler(msgPrinter ifc.MsgPrinter, logPrinter ifc.MsgPrinter, runner ifc.CmdRunner, opts ...*bool) {
	if !hlp.WasRunAsRoot(runner) {
		msgPrinter.Errorfln("avly: flag 'stop' needs to be executed as root")
	}
	vncServerDrained, err := drain(msgPrinter, logPrinter, runner)
	if err != nil {
		logPrinter.Errorfln("avly: %s", err.Error())
	}
	if !vncServerDrained {
		logPrinter.Errorfln("avly: could not drain VNC server")
	}
}

func enterHandler(msgPrinter ifc.MsgPrinter, logPrinter ifc.MsgPrinter, runner ifc.CmdRunner, opts ...*bool) {
	if !hlp.WasRunAsRoot(runner) {
		msgPrinter.Errorfln("avly: flag 'enter' needs to be executed as root")
	}
	logging, wine, _, _, _, err := enter(msgPrinter, logPrinter, runner)
	if err != nil {
		logPrinter.Errorfln("avly: %s", err.Error())
	}
	if !logging {
		logPrinter.Errorfln("avly: could not create logfile")
	}
	if !wine {
		logPrinter.Errorfln("avly: could not install Wine (third-party)")
	}
	logPrinter.Printfln("All set. Watching...")

	var issueMeter, countsUpTo1Day, countsUpTo1Week uint32
	for {
		time.Sleep(60 * time.Second)
		countsUpTo1Day++
		countsUpTo1Week++
		if countsUpTo1Day >= 1440 {
			cleanUpHandler(msgPrinter, logPrinter, runner)
			countsUpTo1Day = 0
		}
		if countsUpTo1Week >= 10080 {
			var locIssueMeter uint32
		COPYTOBAKLOG:
			_, _, errPrToBk := runner.RunCmdSync("cat $AVL_LOGS/avly.log >> $AVL_LOGS/avly.bak-$(date +\"%Y_%m\").log", &env)
			if errPrToBk != nil {
				locIssueMeter++
				if locIssueMeter > 3 {
					logPrinter.Printfln("avly: warn: critical: could not backup logs")
					countsUpTo1Week = 0
					continue
				}
				goto COPYTOBAKLOG
			}
			runner.RunCmdSync("cat /dev/null > $AVL_LOGS/avly.log", &env)
			runner.RunCmdSync("echo $(date +\"%Y/%m/%d %T\") For older logs, see: $(echo $AVL_LOGS/avly.bak-$(date +\"%Y_%m\").log) >> $AVL_LOGS/avly.log", &env)
			countsUpTo1Week = 0
		}
	WATCH:
		pidTarget, _, _ := runner.RunCmdSync("pidof \"terminal64.exe\" | cut -d \" \" -f 1", &env)
		if len(pidTarget) == 0 {
			if issueMeter >= 3 {
				break
			}
			issueMeter++
			logPrinter.Printfln("Force stop before relaunch")
			stopHandler(msgPrinter, logPrinter, runner)
			launchHandler(msgPrinter, logPrinter, runner)
			goto WATCH
		}
		pidX11, _, _ := runner.RunCmdSync("pidof \"x11vnc\" | cut -d \" \" -f 1", &env)
		pidXv, _, _ := runner.RunCmdSync("pidof \"Xvfb\" | cut -d \" \" -f 1", &env)
		if len(pidX11) == 0 || len(pidXv) == 0 {
			if issueMeter >= 3 {
				break
			}
			issueMeter++
			logPrinter.Printfln("Force drain before re-fledge")
			drainHandler(msgPrinter, logPrinter, runner)
			fledgeHandler(msgPrinter, logPrinter, runner)
			goto WATCH
		}
		if issueMeter > 0 {
			issueMeter = 0
			logPrinter.Printfln("All set. Watching...")
		}
	}
}

func prepare(msgPrinter ifc.MsgPrinter, logPrinter ifc.MsgPrinter, runner ifc.CmdRunner) (finishedWineSetup, installedExecutables bool, err error) {
	var dq hlp.ProcDeathQueue
	defer dq.LetDie(runner, &env)

	logPrinter.Printfln("Bee preparation...")

	// STEP 1: Setting up Wine prefix
	err = hlp.PrepareWineprefix(runner, &env)
	if err != nil {
		return
	}
	logPrinter.Printfln("prepare: step 1/2")
	finishedWineSetup = true

	// STEP 2: Install target executable(s)
	_, pIns, errIns := runner.RunCmdAsync("wine $THIRD_PARTY/mt5setup.exe /auto", &env)
	dq.Add(pIns)
	if errIns != nil {
		return
	}
	time.Sleep(45 * time.Second)
	logPrinter.Printfln("prepare: step 2/2")
	installedExecutables = true

	_, pLgR, _ := runner.RunCmdSync("echo $(date +\"%Y/%m/%d %T\") Bee is ready and set >> $AVL_LOGS/avly.log", &env)
	dq.Add(pLgR)
	logPrinter.Printfln("Bee preparation successful")

	return
}

func fledge(msgPrinter ifc.MsgPrinter, logPrinter ifc.MsgPrinter, runner ifc.CmdRunner) (isFrameBufferRunning, isVncServerRunning bool, err error) {
	logPrinter.Printfln("Safely open framebuffer and pull up VNC server...")

	xvfbPid, _, errCmd := runner.RunCmdSync("pidof \"Xvfb\" | cut -d \" \" -f 1", &env)
	if errCmd != nil {
		err = errCmd
		return
	}
	if len(xvfbPid) == 0 {
		logPrinter.Printfln("Framebuffer is not running...")
		runner.RunCmdSync("killall -9 \"i3*\"", &env)
		_, _, errCmd = runner.RunCmdAsync("Xvfb $DISPLAY -screen $SCREEN_NUM $SCREEN_WHD +extension DPMS +extension GLX +extension RANDR +extension RENDER &> $AVL_LOGS/xvfb.log", &env)
		if errCmd != nil {
			err = errCmd
			return
		}
		runner.RunCmdSync("echo $(date +\"%Y/%m/%d %T\") Opened framebuffer >> $AVL_LOGS/avly.log", &env)
	}
	isFrameBufferRunning = true
	logPrinter.Printfln("Framebuffer: OK")

	x11vncPid, _, errCmd := runner.RunCmdSync("pidof \"x11vnc\" | cut -d \" \" -f 1", &env)
	if errCmd != nil {
		err = errCmd
		return
	}
	if len(x11vncPid) == 0 {
		logPrinter.Printfln("VNC server is not running...")
		_, _, errCmd = runner.RunCmdAsync("x11vnc -display $DISPLAY -bg -forever -nopw -quiet -rfbport 5900 -xkb -o $AVL_LOGS/x11vnc.log", &env)
		if errCmd != nil {
			err = errCmd
			return
		}
		runner.RunCmdSync("xset -dpms", &env)
		runner.RunCmdSync("xset s noblank", &env)
		runner.RunCmdSync("xset s off", &env)
		_, _, errCmd = runner.RunCmdAsync("i3 &> $AVL_LOGS/i3.log", &env)
		if errCmd != nil {
			err = errCmd
			return
		}
	}
	isVncServerRunning = true
	runner.RunCmdSync("echo $(date +\"%Y/%m/%d %T\") Pulled up VNC server >> $AVL_LOGS/avly.log", &env)
	logPrinter.Printfln("VNC server: OK")

	return
}

func launch(msgPrinter ifc.MsgPrinter, logPrinter ifc.MsgPrinter, runner ifc.CmdRunner) (isTargetProcessRunning bool, err error) {
	var tcfErr error

	hlp.GetTCF(
		func() {
			// Check for running instances
			targetPid, _ := runner.PanicCmdSync("pidof \"terminal64.exe\" | cut -d \" \" -f 1", &env)
			if len(targetPid) > 0 {
				logPrinter.Printfln("Target process is running")
				isTargetProcessRunning = true
				return
			}

		TARGETRUN:
			// Launch a new instance
			logPrinter.Printfln("Target process is not running...")
			runner.PanicCmdAsync("wine $WINEPREFIX/dosdevices/c\\:/Program\\ Files/MetaTrader\\ 5/terminal64.exe /portable &> $AVL_LOGS/target.log", &env)
			time.Sleep(30 * time.Second)
			targetPid, _ = runner.PanicCmdSync("pidof \"terminal64.exe\" | cut -d \" \" -f 1", &env)
			if len(targetPid) > 0 {
				runner.RunCmdSync("echo $(date +\"%Y/%m/%d %T\") Launched target executable >> $AVL_LOGS/avly.log", &env)
				logPrinter.Printfln("Target process is running")
				isTargetProcessRunning = true
				return
			}
			goto TARGETRUN

		},
		func(caught error) {
			tcfErr = caught
		},
		nil,
	).Run()
	if tcfErr != nil {
		err = tcfErr
	} else {
		logPrinter.Printfln("Target process: OK")
	}

	return
}

func cleanUp(msgPrinter ifc.MsgPrinter, logPrinter ifc.MsgPrinter, runner ifc.CmdRunner) (cleanedUp bool, err error) {
	var dq hlp.ProcDeathQueue
	defer dq.LetDie(runner, &env)

	logPrinter.Printfln("Clean up...")

	var tcfError error
	hlp.GetTCF(
		func() {
			_, pRmLg := runner.PanicCmdSync("rm -rf $PREFIX/dosdevices/c\\:/Program\\ Files/MetaTrader\\ 5/logs/*", &env)
			dq.Add(pRmLg)
			_, pRmHs := runner.PanicCmdSync("rm -rf $PREFIX/dosdevices/c\\:/Program\\ Files/MetaTrader\\ 5/history/*", &env)
			dq.Add(pRmHs)
			_, pRmCs := runner.PanicCmdSync("rm -rf $PREFIX/dosdevices/c\\:/Program\\ Files/MetaTrader\\ 5/*.csv", &env)
			dq.Add(pRmCs)
		},
		func(caught error) {
			tcfError = caught
			_, logProc, _ := runner.RunCmdSync("echo $(date +\"%Y/%m/%d %T\") Problems during cleanup >> $AVL_LOGS/avly.log", &env)
			dq.Add(logProc)
		},
		func() {
			if tcfError == nil {
				_, logProc, _ := runner.RunCmdSync("echo $(date +\"%Y/%m/%d %T\") Cleaned up >> $AVL_LOGS/avly.log", &env)
				dq.Add(logProc)
			}
			dq.LetDie(runner, &env)
		},
	).Run()
	err = tcfError

	return
}

func stop(msgPrinter ifc.MsgPrinter, logPrinter ifc.MsgPrinter, runner ifc.CmdRunner) (targetProcessDead bool, err error) {
	var dq hlp.ProcDeathQueue
	defer dq.LetDie(runner, &env)

	logPrinter.Printfln("Stop target process(es)...")

	var targetProcessNames, targetPids []string
	targetProcessNames = append(targetProcessNames, "terminal64.exe")
	for {
		targetPids = nil
		for targetProcCount := 0; targetProcCount < len(targetProcessNames); targetProcCount++ {
			pid, pidOfProc, _ := runner.RunCmdSync(fmt.Sprintf("pidof \"%s\" | cut -d \" \" -f 1", targetProcessNames[targetProcCount]), &env)
			dq.Add(pidOfProc)
			if len(pid) > 0 {
				targetPids = append(targetPids, pid)
			}
		}
		if targetPids == nil {
			targetProcessDead = true
			break
		}
		for targetPidCount := 0; targetPidCount < len(targetPids); targetPidCount++ {
			_, pkillProc, _ := runner.RunCmdSync(fmt.Sprintf("pkill -15 \"%s\"", targetPids[targetPidCount]), &env)
			dq.Add(pkillProc)
		}
	}
	_, logProc, _ := runner.RunCmdSync("echo $(date +\"%Y/%m/%d %T\") Stopped target process(es) >> $AVL_LOGS/avly.log", &env)
	dq.Add(logProc)
	logPrinter.Printfln("Stopped target process(es)")

	return
}

func drain(msgPrinter ifc.MsgPrinter, logPrinter ifc.MsgPrinter, runner ifc.CmdRunner) (vncServerDrained bool, err error) {
	var dq hlp.ProcDeathQueue
	defer dq.LetDie(runner, &env)

	logPrinter.Printfln("Drain VNC server...")

	var xProcNames, targetPids []string
	xProcNames = append(xProcNames, "x11vnc", "Xvfb")
	for {
		targetPids = nil
		for targetProcCount := 0; targetProcCount < len(xProcNames); targetProcCount++ {
			pid, pidOfProc, _ := runner.RunCmdSync(fmt.Sprintf("pidof \"%s\" | cut -d \" \" -f 1", xProcNames[targetProcCount]), &env)
			dq.Add(pidOfProc)
			if len(pid) > 0 {
				targetPids = append(targetPids, pid)
			}
		}
		if targetPids == nil {
			vncServerDrained = true
			break
		}
		for targetPidCount := 0; targetPidCount < len(targetPids); targetPidCount++ {
			_, pkillProc, _ := runner.RunCmdSync(fmt.Sprintf("pkill -9 \"%s\"", targetPids[targetPidCount]), &env)
			dq.Add(pkillProc)
		}
	}
	_, logProc, _ := runner.RunCmdSync("echo $(date +\"%Y/%m/%d %T\") Drained VNC server >> $AVL_LOGS/avly.log", &env)
	dq.Add(logProc)
	logPrinter.Printfln("Drained VNC server")

	return
}

func enter(msgPrinter ifc.MsgPrinter, logPrinter ifc.MsgPrinter, runner ifc.CmdRunner) (enabledLogging, installedWine, isFledged, isPrepared, isLaunched bool, err error) {
	var dq hlp.ProcDeathQueue
	defer dq.LetDie(runner, &env)

	logPrinter.Printfln("Start initialization...")

	// STEP 1: Prepare logging
	_, pPrLg, errAvLog := runner.RunCmdSync("cat /dev/null > $AVL_LOGS/avly.log", &env)
	dq.Add(pPrLg)
	if errAvLog != nil {
		err = errAvLog
		return
	}
	logPrinter.Printfln("enter: step 1/5")
	enabledLogging = true

	// STEP 2: Install Wine
	errWineInstall := hlp.InstallWine(runner, &env)
	if errWineInstall != nil {
		return
	}
	logPrinter.Printfln("enter: step 2/5")
	installedWine = true

	// STEP 3: Fledge VNC server
	fledgeHandler(msgPrinter, logPrinter, runner)
	logPrinter.Printfln("enter: step 3/5")
	isFledged = true

	// STEP 4: Prepare bee
	prepareHandler(msgPrinter, logPrinter, runner)
	logPrinter.Printfln("enter: step 4/5")
	isPrepared = true

	// STEP 5: Initial target launch
	launchHandler(msgPrinter, logPrinter, runner)
	logPrinter.Printfln("enter: step 5/5")
	isLaunched = true

	_, pLgW, _ := runner.RunCmdSync("echo $(date +\"%Y/%m/%d %T\") Bee is now working >> $AVL_LOGS/avly.log", &env)
	dq.Add(pLgW)
	logPrinter.Printfln("Initialization successful")

	return
}
