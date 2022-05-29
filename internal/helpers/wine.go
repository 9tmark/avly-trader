// Copyright (C) 2022 The Avly Trader Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package helpers

import (
	"time"

	ifc "github.com/9tmark/avly-trader/internal/interfaces"
)

func InstallWine(runner ifc.CmdRunner, env *[]string) (err error) {
	var dq ProcDeathQueue
	GetTCF(
		func() {
			_, pWg1 := runner.PanicCmdSync("wget -nc https://dl.winehq.org/wine-builds/winehq.key -P /usr/share/keyrings", env)
			dq.Add(pWg1)
			_, pMv := runner.PanicCmdSync("mv /usr/share/keyrings/winehq.key /usr/share/keyrings/winehq-archive.key", env)
			dq.Add(pMv)
			_, pWg2 := runner.PanicCmdSync("wget -nc https://dl.winehq.org/wine-builds/ubuntu/dists/focal/winehq-focal.sources -P /etc/apt/sources.list.d", env)
			dq.Add(pWg2)
			_, pApU := runner.PanicCmdSync("apt-get update -yq", env)
			dq.Add(pApU)
			_, pApI := runner.PanicCmdSync("apt-get install -yq --install-recommends winehq-staging=7.2~focal-1 wine-staging=7.2~focal-1 wine-staging-amd64=7.2~focal-1 wine-staging-i386=7.2~focal-1", env)
			dq.Add(pApI)
		},
		func(caught error) {
			err = caught
		},
		func() {
			dq.LetDie(runner, env)
		},
	).Run()

	return
}

func PrepareWineprefix(runner ifc.CmdRunner, env *[]string) (err error) {
	var dq ProcDeathQueue
	GetTCF(
		func() {
			_, pWb1 := runner.PanicCmdAsync("wine wineboot -u &> $AVL_LOGS/wine.log", env)
			dq.Add(pWb1)
			time.Sleep(20 * time.Second)
			_, pRet := runner.PanicCmdAsync("xdotool key --clearmodifiers Return", env)
			dq.Add(pRet)
			time.Sleep(80 * time.Second)
			runner.PanicCmdAsync("wine wineboot -u &>> $AVL_LOGS/wine.log", env)
			time.Sleep(20 * time.Second)
			_, pMon := runner.PanicCmdAsync("wine msiexec /i $THIRD_PARTY/wine-mono-7.1.1-x86.msi &>> $AVL_LOGS/wine.log", env)
			dq.Add(pMon)
			time.Sleep(45 * time.Second)
			_, pGec := runner.PanicCmdAsync("wine msiexec /i $THIRD_PARTY/wine_gecko-2.47-x86_64.msi &>> $AVL_LOGS/wine.log", env)
			dq.Add(pGec)
			time.Sleep(20 * time.Second)
			_, pCp := runner.PanicCmdSync("cp $THIRD_PARTY/winetricks /usr/local/bin", env)
			dq.Add(pCp)
			_, pCh := runner.PanicCmdSync("chmod +x /usr/local/bin/winetricks", env)
			dq.Add(pCh)
			_, pWt := runner.PanicCmdAsync("winetricks -f --unattended corefonts", env)
			dq.Add(pWt)
			time.Sleep(45 * time.Second)
		},
		func(caught error) {
			err = caught
		},
		func() {
			dq.LetDie(runner, env)
		},
	).Run()

	return
}
