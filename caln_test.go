package caln_test

import (
	"caln"
	"testing"
	"time"
)

// TestBSDateFromGregorian checks conversions against published Bikram Sambat
// calendars and independent datasets. Expected values are taken from:
//
//   - Nepal Panchanga Nirnayak Samiti / Hamro Patro and KhataIN calendars
//     (used for recent years such as 2081-2083 BS).
//   - The Nepali Calendar (nepcal.com) BS<->AD export, BS 1975-2100:
//     https://github.com/mesaugat/bikram-sambat-anno-domini-fixtures
//   - The ashesh.com.np dataset used by the Temporal Nepali calendar cookbook:
//     https://tc39.es/proposal-temporal/docs/cookbook-nepali-calendar.html
//   - Wikipedia's Timeline of Nepalese politics for the historical mappings.
func TestBSDateFromGregorian(t *testing.T) {
	ad := func(year int, month time.Month, day int) time.Time {
		return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	}
	bs := func(year int, month caln.BSMonth, day int) caln.BSDate {
		return caln.BSDate{Year: year, Month: month, Day: day}
	}

	tests := []struct {
		name      string
		gregorian time.Time
		want      caln.BSDate
	}{
		// Era anchors and the start of Baisakh across the supported range.
		{name: "1943-04-14 AD is 2000 Baisakh 1 BS", gregorian: ad(1943, time.April, 14), want: bs(2000, caln.Baisakh, 1)},
		{name: "1953-04-13 AD is 2010 Baisakh 1 BS", gregorian: ad(1953, time.April, 13), want: bs(2010, caln.Baisakh, 1)},
		{name: "1963-04-14 AD is 2020 Baisakh 1 BS", gregorian: ad(1963, time.April, 14), want: bs(2020, caln.Baisakh, 1)},
		{name: "1983-04-14 AD is 2040 Baisakh 1 BS", gregorian: ad(1983, time.April, 14), want: bs(2040, caln.Baisakh, 1)},
		{name: "2003-04-14 AD is 2060 Baisakh 1 BS", gregorian: ad(2003, time.April, 14), want: bs(2060, caln.Baisakh, 1)},
		{name: "2013-04-14 AD is 2070 Baisakh 1 BS", gregorian: ad(2013, time.April, 14), want: bs(2070, caln.Baisakh, 1)},
		{name: "2023-04-14 AD is 2080 Baisakh 1 BS", gregorian: ad(2023, time.April, 14), want: bs(2080, caln.Baisakh, 1)},
		{name: "2033-04-14 AD is 2090 Baisakh 1 BS", gregorian: ad(2033, time.April, 14), want: bs(2090, caln.Baisakh, 1)},
		// Gregorian year boundaries and leap days.
		{name: "1999-12-31 AD is 2056 Paush 16 BS", gregorian: ad(1999, time.December, 31), want: bs(2056, caln.Paush, 16)},
		{name: "2000-01-01 AD is 2056 Paush 17 BS", gregorian: ad(2000, time.January, 1), want: bs(2056, caln.Paush, 17)},
		{name: "2000-02-29 AD is 2056 Falgun 17 BS", gregorian: ad(2000, time.February, 29), want: bs(2056, caln.Falgun, 17)},
		{name: "2024-02-29 AD is 2080 Falgun 17 BS", gregorian: ad(2024, time.February, 29), want: bs(2080, caln.Falgun, 17)},
		// Recent and historical dates confirmed by Nepali calendar sources.
		{name: "2000-04-13 AD is 2057 Baisakh 1 BS", gregorian: ad(2000, time.April, 13), want: bs(2057, caln.Baisakh, 1)},
		{name: "2015-04-14 AD is 2072 Baisakh 1 BS", gregorian: ad(2015, time.April, 14), want: bs(2072, caln.Baisakh, 1)},
		{name: "2015-04-25 AD is 2072 Baisakh 12 BS", gregorian: ad(2015, time.April, 25), want: bs(2072, caln.Baisakh, 12)},
		{name: "2008-05-28 AD is 2065 Jestha 15 BS", gregorian: ad(2008, time.May, 28), want: bs(2065, caln.Jestha, 15)},
		{name: "2025-04-14 AD is 2082 Baisakh 1 BS", gregorian: ad(2025, time.April, 14), want: bs(2082, caln.Baisakh, 1)},
		{name: "2026-04-13 AD is 2082 Chaitra 30 BS", gregorian: ad(2026, time.April, 13), want: bs(2082, caln.Chaitra, 30)},
		{name: "2026-04-14 AD is 2083 Baisakh 1 BS", gregorian: ad(2026, time.April, 14), want: bs(2083, caln.Baisakh, 1)},
		{name: "2024-06-14 AD is 2081 Jestha 32 BS", gregorian: ad(2024, time.June, 14), want: bs(2081, caln.Jestha, 32)},
		{name: "2024-06-15 AD is 2081 Ashadh 1 BS", gregorian: ad(2024, time.June, 15), want: bs(2081, caln.Ashadh, 1)},
		{name: "2022-04-14 AD is 2079 Baisakh 1 BS", gregorian: ad(2022, time.April, 14), want: bs(2079, caln.Baisakh, 1)},
		{name: "2022-05-15 AD is 2079 Jestha 1 BS", gregorian: ad(2022, time.May, 15), want: bs(2079, caln.Jestha, 1)},
		{name: "2022-06-15 AD is 2079 Ashadh 1 BS", gregorian: ad(2022, time.June, 15), want: bs(2079, caln.Ashadh, 1)},
		{name: "2022-07-17 AD is 2079 Shrawan 1 BS", gregorian: ad(2022, time.July, 17), want: bs(2079, caln.Shrawan, 1)},
		{name: "2022-08-17 AD is 2079 Bhadra 1 BS", gregorian: ad(2022, time.August, 17), want: bs(2079, caln.Bhadra, 1)},
		{name: "2022-09-17 AD is 2079 Ashoj 1 BS", gregorian: ad(2022, time.September, 17), want: bs(2079, caln.Ashoj, 1)},
		{name: "2022-10-18 AD is 2079 Kartik 1 BS", gregorian: ad(2022, time.October, 18), want: bs(2079, caln.Kartik, 1)},
		{name: "2022-11-17 AD is 2079 Mangsir 1 BS", gregorian: ad(2022, time.November, 17), want: bs(2079, caln.Mangsir, 1)},
		{name: "2022-12-16 AD is 2079 Paush 1 BS", gregorian: ad(2022, time.December, 16), want: bs(2079, caln.Paush, 1)},
		{name: "2023-01-15 AD is 2079 Magh 1 BS", gregorian: ad(2023, time.January, 15), want: bs(2079, caln.Magh, 1)},
		{name: "2023-02-13 AD is 2079 Falgun 1 BS", gregorian: ad(2023, time.February, 13), want: bs(2079, caln.Falgun, 1)},
		{name: "2023-03-15 AD is 2079 Chaitra 1 BS", gregorian: ad(2023, time.March, 15), want: bs(2079, caln.Chaitra, 1)},
		{name: "2022-05-14 AD is 2079 Baisakh 31 BS", gregorian: ad(2022, time.May, 14), want: bs(2079, caln.Baisakh, 31)},
		{name: "2022-07-16 AD is 2079 Ashadh 32 BS", gregorian: ad(2022, time.July, 16), want: bs(2079, caln.Ashadh, 32)},
		{name: "2022-12-15 AD is 2079 Mangsir 29 BS", gregorian: ad(2022, time.December, 15), want: bs(2079, caln.Mangsir, 29)},
		{name: "2023-04-13 AD is 2079 Chaitra 30 BS", gregorian: ad(2023, time.April, 13), want: bs(2079, caln.Chaitra, 30)},
		// Historical mappings from Wikipedia's Timeline of Nepalese politics.
		{name: "1991-05-12 AD is 2048 Baisakh 29 BS", gregorian: ad(1991, time.May, 12), want: bs(2048, caln.Baisakh, 29)},
		{name: "1991-05-29 AD is 2048 Jestha 15 BS", gregorian: ad(1991, time.May, 29), want: bs(2048, caln.Jestha, 15)},
		{name: "1994-07-11 AD is 2051 Ashadh 27 BS", gregorian: ad(1994, time.July, 11), want: bs(2051, caln.Ashadh, 27)},
		{name: "1995-06-13 AD is 2052 Jestha 30 BS", gregorian: ad(1995, time.June, 13), want: bs(2052, caln.Jestha, 30)},
		{name: "1995-08-28 AD is 2052 Bhadra 12 BS", gregorian: ad(1995, time.August, 28), want: bs(2052, caln.Bhadra, 12)},
		{name: "1996-02-13 AD is 2052 Falgun 1 BS", gregorian: ad(1996, time.February, 13), want: bs(2052, caln.Falgun, 1)},
		{name: "1997-09-17 AD is 2054 Ashoj 1 BS", gregorian: ad(1997, time.September, 17), want: bs(2054, caln.Ashoj, 1)},
		{name: "1998-03-05 AD is 2054 Falgun 21 BS", gregorian: ad(1998, time.March, 5), want: bs(2054, caln.Falgun, 21)},
		{name: "1999-04-19 AD is 2056 Baisakh 6 BS", gregorian: ad(1999, time.April, 19), want: bs(2056, caln.Baisakh, 6)},
		{name: "1999-05-31 AD is 2056 Jestha 17 BS", gregorian: ad(1999, time.May, 31), want: bs(2056, caln.Jestha, 17)},
		{name: "2000-02-16 AD is 2056 Falgun 4 BS", gregorian: ad(2000, time.February, 16), want: bs(2056, caln.Falgun, 4)},
		{name: "2001-07-23 AD is 2058 Shrawan 8 BS", gregorian: ad(2001, time.July, 23), want: bs(2058, caln.Shrawan, 8)},
		{name: "2002-05-22 AD is 2059 Jestha 8 BS", gregorian: ad(2002, time.May, 22), want: bs(2059, caln.Jestha, 8)},
		{name: "2003-06-05 AD is 2060 Jestha 22 BS", gregorian: ad(2003, time.June, 5), want: bs(2060, caln.Jestha, 22)},
		{name: "2004-06-04 AD is 2061 Jestha 22 BS", gregorian: ad(2004, time.June, 4), want: bs(2061, caln.Jestha, 22)},
		{name: "2005-02-01 AD is 2061 Magh 19 BS", gregorian: ad(2005, time.February, 1), want: bs(2061, caln.Magh, 19)},
		{name: "2006-04-24 AD is 2063 Baisakh 11 BS", gregorian: ad(2006, time.April, 24), want: bs(2063, caln.Baisakh, 11)},
		{name: "2007-01-15 AD is 2063 Magh 1 BS", gregorian: ad(2007, time.January, 15), want: bs(2063, caln.Magh, 1)},
		{name: "2008-04-10 AD is 2064 Chaitra 28 BS", gregorian: ad(2008, time.April, 10), want: bs(2064, caln.Chaitra, 28)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := caln.BSDateFromGregorian(tt.gregorian)
			if got != tt.want {
				t.Errorf(
					"BSDateFromGregorian(%s) = %+v, want %+v",
					tt.gregorian.Format("2006-01-02"),
					got,
					tt.want,
				)
			}
		})
	}
}
