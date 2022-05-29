// Copyright (C) 2022 The Avly Trader Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package helpers

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	ifc "github.com/9tmark/avly-trader/internal/interfaces"
)

func TestRunsBlocksInOrderNoCatch(t *testing.T) {
	var pSpy ifc.SpyMsgPrinter
	first, second := "I am", "second"

	GetTCF(
		func() {
			pSpy.Printfln(first)
		},
		func(ignored error) {},
		func() {
			pSpy.Printfln(second)
		},
	).Run()

	calls, lastMsg := pSpy.Calls, pSpy.LastMessage
	if calls != 2 {
		t.Errorf("calls: Expected '%d' to be '%d'", calls, 2)
	}
	if strings.Compare(lastMsg, second) != 0 {
		t.Errorf("lastMsg: Expected '%s' to be '%s'", lastMsg, second)
	}
}

func TestRunsBlocksInOrderCatchingUnkownProblem(t *testing.T) {
	var pSpy ifc.SpyMsgPrinter
	var caughtMsg string
	first, second, third := "I", "am", "third"

	GetTCF(
		func() {
			pSpy.Printfln(first)
			panic(42)
		},
		func(caught error) {
			caughtMsg = caught.Error()
			pSpy.Printfln(second)
		},
		func() {
			pSpy.Printfln(third)
		},
	).Run()

	calls, history := pSpy.Calls, pSpy.History
	if calls != 3 {
		t.Errorf("calls: Expected '%d' to be '%d'", calls, 3)
	}
	if strings.Compare(history[1], second) != 0 {
		t.Errorf("history[1]: Expected '%s' to be '%s'", history[1], second)
	}
	if strings.Compare(history[2], third) != 0 {
		t.Errorf("history[2]: Expected '%s' to be '%s'", history[2], third)
	}
	if strings.Compare(caughtMsg, "paniced for unknown problem") != 0 {
		t.Errorf("caughtMsg: Expected '%s' to be '%s'", caughtMsg, "paniced for unknown problem")
	}
}

func TestRunsBlocksInOrderCatchingError(t *testing.T) {
	var finalEvidence, caughtMsg string

	GetTCF(
		func() {
			panic(errors.New("hammertime error: the trouser's legs are far too wide"))
		},
		func(howeverCaughtAs error) {
			caughtMsg = howeverCaughtAs.Error()
		},
		func() {
			finalEvidence = "worked out"
		},
	).Run()

	if strings.Compare(finalEvidence, "worked out") != 0 {
		t.Errorf("finalEvidence: Expected '%s' to be '%s'", finalEvidence, "worked out")
	}
	if !strings.HasPrefix(caughtMsg, "hammertime error") {
		t.Errorf("caughtMsg: Expected '%v' to be a 'hammertime error'", caughtMsg)
	}
}

func TestRunsBlocksInOrderCatchingString(t *testing.T) {
	var finalEvidence, caughtMsg string

	GetTCF(
		func() {
			panic("not an error")
		},
		func(howeverCaughtAs error) {
			caughtMsg = howeverCaughtAs.Error()
		},
		func() {
			finalEvidence = "worked out"
		},
	).Run()

	if strings.Compare(finalEvidence, "worked out") != 0 {
		t.Errorf("finalEvidence: Expected '%s' to be '%s'", finalEvidence, "worked out")
	}
	if strings.Compare(caughtMsg, "not an error") != 0 {
		t.Errorf("caughtMsg: Expected '%s' to be '%s'", caughtMsg, "not an error")
	}
}

func TestPanicsForNilFunctions(t *testing.T) {
	defer func() {
		if panicString := fmt.Sprintf("%v", recover()); !strings.HasPrefix(panicString, "catchgroup error") {
			t.Errorf("panicString: Expected '%v' to be a 'catchgroup error'", panicString)
		}
	}()

	GetTCF(nil, nil, nil).Run()
}
