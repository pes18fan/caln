/*
This code is adapted from github.com/starryalley/go-julianday, which is
published under the following license:

MIT License

Copyright (c) 2020 Mark Kuo

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package julianday

import (
	"time"
)

const (
	secondsInADay       = 86400
	unixEpochJulianDate = 2440587.5
)

// Number (JDN) returns the integer of Julian Day on the noon (12:00PM)
// of this date.
func Number(Y int, M time.Month, D int) int64 {
	// https://en.wikipedia.org/wiki/Julian_day#Julian_day_number_calculation
	return (1461*(int64(Y)+4800+(int64(M)-14)/12))/4 +
		(367*(int64(M)-2-12*((int64(M)-14)/12)))/12 -
		(3*((int64(Y)+4900+(int64(M)-14)/12)/100))/4 +
		int64(D) - 32075
}

// Date returns the Julian date of time
func Date(t time.Time) float64 {
	return float64(t.Unix())/secondsInADay + unixEpochJulianDate
}

// ToTime returns a time.Time for the Julian Day instant jd
func ToTime(jd float64) time.Time {
	return time.Unix(int64((jd-unixEpochJulianDate)*secondsInADay), 0)
}
