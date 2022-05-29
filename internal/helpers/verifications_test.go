// Copyright (C) 2022 The Avly Trader Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package helpers

import (
	"testing"

	ifc "github.com/9tmark/avly-trader/internal/interfaces"
)

func TestHasOnlyOneTrueValueSingleTrueValue(t *testing.T) {
	b1 := true

	if result := HasOnlyOneTrueValue(&b1); !result {
		t.Errorf("Expected '%t' to be '%t'", result, true)
	}
}

func TestHasOnlyOneTrueValueSingleFalseValue(t *testing.T) {
	b1 := false

	if result := HasOnlyOneTrueValue(&b1); result {
		t.Errorf("Expected '%t' to be '%t'", result, false)
	}
}

func TestHasOnlyOneTrueValueMixedValuesPositive(t *testing.T) {
	b1, b2, b3 := false, true, false

	if result := HasOnlyOneTrueValue(&b1, &b2, &b3); !result {
		t.Errorf("Expected '%t' to be '%t'", result, true)
	}
}

func TestHasOnlyOneTrueValueMixedValuesNegative(t *testing.T) {
	b1, b2, b3 := false, true, true

	if result := HasOnlyOneTrueValue(&b1, &b2, &b3); result {
		t.Errorf("Expected '%t' to be '%t'", result, false)
	}
}

func TestHasOnlyOneTrueValueAllTrueValues(t *testing.T) {
	b1, b2, b3 := true, true, true

	if result := HasOnlyOneTrueValue(&b1, &b2, &b3); result {
		t.Errorf("Expected '%t' to be '%t'", result, false)
	}
}

func TestHasOnlyOneTrueValueAllFalseValues(t *testing.T) {
	b1, b2, b3 := false, false, false

	if result := HasOnlyOneTrueValue(&b1, &b2, &b3); result {
		t.Errorf("Expected '%t' to be '%t'", result, false)
	}
}

func TestWasRunAsRootNegative(t *testing.T) {
	/* As the execution function of the spy simply handovers the given command
	(here: 'whoami'), the result should be false for fictional user 'whoami' not being root
	*/
	if result := WasRunAsRoot(&ifc.SpySafeCmdRunner{}); result {
		t.Errorf("Expected '%t' to be '%t'", result, false)
	}
}
