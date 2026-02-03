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

package datatype

import (
	"testing"
	"time"
)

//=============================================================================

func TestToIntDate(t *testing.T) {
	td := time.Date(2025, 05, 03, 11, 12, 13, 0, time.UTC)
	id := ToIntDate(&td)
	exp := IntDate(20250503)

	if id != exp {
		t.Errorf("ToIntDate failed. Expected %v but got %v", exp, id)
	}
}

//=============================================================================

func TestString(t *testing.T) {
	id := IntDate(20250503)
	sd := id.String()
	exp := "2025-05-03"

	if sd != exp {
		t.Errorf("ToDate failed. Expected %v but got %v", exp, sd)
	}
}

//=============================================================================

func TestParseIntDate(t *testing.T) {
	sd := "20250503"
	exp := IntDate(20250503)

	id, err := ParseIntDate(sd, true)
	if err != nil {
		t.Errorf("ParseIntDate failed. Expected %v but got %v", exp, id)
	}

	//---

	sd = "-20250503"

	id, err = ParseIntDate(sd, true)
	if err == nil {
		t.Errorf("ParseIntDate failed. Date is indicated as valid but it is not: %v", id)
	}

	//---

	sd = ""

	id, err = ParseIntDate(sd, false)
	if err != nil || !id.IsNil() {
		t.Errorf("ParseIntDate failed. Date is nil but got a valid date: %v", id)
	}
}

//=============================================================================

func TestDays(t *testing.T) {
	s := IntDate(20250503)
	d := IntDate(20250505)

	if s.Days(d) != 2 {
		t.Errorf("Days failed. Expected %v but got %v", 2, s.Days(d))
	}
}

//=============================================================================

func TestDaysLeap(t *testing.T) {
	s := IntDate(20240302)
	d := IntDate(20240228)

	if s.Days(d) != -3 {
		t.Errorf("Days failed. Expected %v but got %v", -3, s.Days(d))
	}
}

//=============================================================================
