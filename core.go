//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package core

import (
	"log/slog"
	"os"
)

//=============================================================================

type Application struct {
	BindAddress string
	Production  bool
	Debug       bool
}

//=============================================================================

type Database struct {
	Address  string
	Name     string
	Username string
	Password string
}

//=============================================================================

type Authentication struct {
	Authority    string
	ClientId     string
	ClientSecret string
}

//=============================================================================

type Platform struct {
	System    string
	Inventory string
	Data      string
	Storage   string
	Portfolio string
}

//=============================================================================

type Journal struct {
	Directory       string
	QueueSize       int
	CompactMessages int
	DbSpoolInterval int
}

//=============================================================================

type Messaging struct {
	Address    string
	Username   string
	Password   string
	Journal    *Journal
}

//=============================================================================

func ExitIfError(err error) {
	if err != nil {
		ExitWithMessage(err.Error())
	}
}

//=============================================================================

func ExitWithMessage(message string) {
	slog.Error(message)
	os.Exit(1)
}

//=============================================================================
