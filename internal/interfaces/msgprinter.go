// Copyright (C) 2022 The Avly Trader Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package interfaces

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type MsgPrinter interface {
	Printfln(msg string, a ...any)
	Errorfln(msg string, a ...any)
}

type SpyMsgPrinter struct {
	Calls       int
	History     []string
	LastMessage string
	LastError   string
}

type FmtMsgPrinter struct{}

type LogMsgPrinter struct{}

func (s *SpyMsgPrinter) Printfln(msg string, a ...any) {
	s.Calls++
	s.History = append(s.History, fmt.Sprintf(msg, a...))
	s.LastMessage = fmt.Sprintf(msg, a...)
}

func (f *FmtMsgPrinter) Printfln(msg string, a ...any) {
	b := strings.Builder{}
	b.WriteString(msg)
	b.WriteString("\n")
	_, _ = fmt.Printf(b.String(), a...)
}

func (l *LogMsgPrinter) Printfln(msg string, a ...any) {
	b := strings.Builder{}
	b.WriteString(msg)
	b.WriteString("\n")
	log.Printf(b.String(), a...)
}

func (s *SpyMsgPrinter) Errorfln(msg string, a ...any) {
	s.Calls++
	s.History = append(s.History, fmt.Sprintf(msg, a...))
	s.LastError = fmt.Sprintf(msg, a...)
}

func (f *FmtMsgPrinter) Errorfln(msg string, a ...any) {
	b := strings.Builder{}
	b.WriteString(msg)
	b.WriteString("\n")
	os.Stderr.WriteString(fmt.Errorf(b.String(), a...).Error())
	os.Exit(1)
}

func (l *LogMsgPrinter) Errorfln(msg string, a ...any) {
	b := strings.Builder{}
	b.WriteString(msg)
	b.WriteString("\n")
	log.New(os.Stderr, "", 0).Fatalf(b.String(), a...)
}
