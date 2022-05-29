// Copyright (C) 2022 The Avly Trader Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package helpers

import (
	"errors"
)

// CatchGroup enables generic error handling for a block of statements which will run through unless one statement throws an error. This way there is no need for a redundant `if err != nil {...}` every time.
// try and catch are mandatory.
type CatchGroup struct {
	try     func()
	catch   func(error)
	finally func()
}

func GetTCF(tryFunc func(), catchFunc func(error), finallyFunc func()) CatchGroup {
	if tryFunc == nil || catchFunc == nil {
		panic(errors.New("catchgroup error: try and catch cannot be nil"))
	}

	return CatchGroup{try: tryFunc, catch: catchFunc, finally: finallyFunc}
}

func (cg CatchGroup) Run() {
	if cg.finally != nil {
		defer cg.finally()
	}
	defer func() {
		if routineErr := recover(); routineErr != nil {
			switch actual := routineErr.(type) {
			case error:
				cg.catch(actual)
			case string:
				cg.catch(errors.New(actual))
			default:
				unknownReason := errors.New("paniced for unknown problem")
				cg.catch(unknownReason)
			}
		}
	}()
	cg.try()
}
