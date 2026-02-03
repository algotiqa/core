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
	"errors"
	"fmt"
	"strconv"
	"time"
)

//=============================================================================

type IntDate int

//=============================================================================

func (dt IntDate) Year() int {
	return int(dt / 10000)
}

//=============================================================================

func (dt IntDate) Month() int {
	return int((dt / 100) % 100)
}

//=============================================================================

func (dt IntDate) Day() int {
	return int(dt % 100)
}

//=============================================================================

func (dt IntDate) String() string {
	return fmt.Sprintf("%4d-%02d-%02d", dt.Year(), dt.Month(), dt.Day())
}

//=============================================================================

func (dt IntDate) IsNil() bool {
	return dt == 0
}

//=============================================================================

func (dt IntDate) IsValid() bool {
	if dt < 0 {
		return false
	}

	d := dt.Day()
	m := dt.Month()

	if m < 1 || m > 12 {
		return false
	}

	if m == 4 || m == 6 || m == 9 || m == 11 {
		return d >= 1 && d <= 30
	}

	if m == 2 {
		return d >= 1 && d <= 29
	}

	return d >= 1 && d <= 31
}

//=============================================================================

func (dt IntDate) ToDateTime(endDay bool, loc *time.Location) time.Time {
	hh := 0
	mm := 0
	ss := 0

	if endDay {
		hh = 23
		mm = 59
		ss = 59
	}

	return time.Date(dt.Year(), time.Month(dt.Month()), dt.Day(), hh, mm, ss, 0, loc)
}

//=============================================================================

func (dt IntDate) AddDays(days int) IntDate {
	t := dt.ToDateTime(false, time.UTC)
	t = t.Add(time.Duration(days) * 24 * time.Hour)

	y, m, d := t.Date()

	return IntDate(y*10000 + int(m)*100 + d)
}

//=============================================================================

func (dt IntDate) IsToday(loc *time.Location) bool {
	return dt == Today(loc)
}

//=============================================================================

func (dt IntDate) Days(d IntDate) int {
	ts := dt.ToDateTime(false, time.UTC)
	td := d.ToDateTime(false, time.UTC)
	seconds := td.Unix() - ts.Unix()

	return int(seconds / 3600 / 24)
}

//=============================================================================
//===
//=== General functions
//===
//=============================================================================

func ToIntDate(t *time.Time) IntDate {
	if t == nil {
		return 0
	}

	y, m, d := t.Date()

	return IntDate(y*10000 + int(m)*100 + d)
}

//=============================================================================

func ParseIntDate(value string, required bool) (IntDate, error) {
	if value == "" {
		if required {
			return 0, errors.New("Value is required")
		}

		return 0, nil
	}

	d, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	id := IntDate(d)

	if !id.IsValid() {
		return 0, errors.New("Invalid date")
	}

	return id, nil
}

//=============================================================================

func Today(loc *time.Location) IntDate {
	tn := time.Now().In(loc)
	return ToIntDate(&tn)
}

//=============================================================================
