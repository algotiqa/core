//=============================================================================
/*
Copyright © 2026 Andrea Carboni andrea.carboni71@gmail.com

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
	"errors"
	"fmt"
	"strconv"
)

//=============================================================================

const NilValue int = 9999

type IntTime int

//=============================================================================

func (it IntTime) Hour() int {
	return int(it / 100)
}

//=============================================================================

func (it IntTime) Minute() int {
	return int(it % 100)
}

//=============================================================================

func (it IntTime) String() string {
	return fmt.Sprintf("%02d:%02d", it.Hour(), it.Minute())
}

//=============================================================================

func (it IntTime) IsNil() bool {
	return it == 9999
}

//=============================================================================

func (it IntTime) IsValid() bool {
	if it < 0 {
		return false
	}

	h := it.Hour()
	m := it.Minute()

	if h < 0 || h > 23 {
		return false
	}
	if m < 0 || m > 59 {
		return false
	}

	return true
}

//=============================================================================

func (it IntTime) AddMinutes(mins int) IntTime {
	totMins := it.Hour()*60 + it.Minute()
	finMins := totMins + mins

	hh := finMins / 60
	mm := finMins % 60

	return NewIntTime(hh, mm)
}

//=============================================================================
//===
//=== General functions
//===
//=============================================================================

func NewIntTime(hours, mins int) IntTime {
	return IntTime(hours*100 + mins)
}

//=============================================================================

func ParseIntTime(value string, required bool) (IntTime, error) {
	if value == "" {
		if required {
			return 0, errors.New("Value is required")
		}

		return IntTime(NilValue), nil
	}

	t, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	it := IntTime(t)

	if !it.IsValid() {
		return 0, errors.New("Invalid time")
	}

	return it, nil
}

//=============================================================================
