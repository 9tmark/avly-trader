// Copyright (C) 2022 The Avly Trader Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package helpers

import (
	"strings"

	ifc "github.com/9tmark/avly-trader/internal/interfaces"
)

var SafeLinuxEnv = []string{
	"PATH=/usr/local/bin:/usr/bin:/usr/local/sbin",
}

func HasOnlyOneTrueValue(exps ...*bool) (eval bool) {
	var foundFirstTrue bool
	for i := 0; i < len(exps); i++ {
		if *exps[i] {
			if !foundFirstTrue {
				eval = true
				foundFirstTrue = true
				continue
			}
			eval = false
		}
	}

	return
}

func IsAvailableInEnvironment(execWord string, r ifc.CmdRunner, env *[]string) (eval bool) {
	out, _, _ := r.RunCmdSync(strings.Join([]string{"which", execWord}, " "), env)
	if !strings.HasSuffix(out, "not found") {
		eval = true
	}

	return
}

func WasRunAsRoot(r ifc.CmdRunner) (eval bool) {
	out, _, _ := r.RunCmdSync("whoami", &SafeLinuxEnv)
	if strings.TrimSpace(out) == "root" {
		eval = true
	}

	return
}
