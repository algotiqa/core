//=============================================================================
/*
Copyright © 2023 Andrea Carboni andrea.carboni71@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
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
