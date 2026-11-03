package caln

import (
	"math"
	"time"
)

func init() {
	initVsopTable()
}

var longBSMonthNames = []string{
	"Baisakh",
	"Jestha",
	"Ashadh",
	"Shrawan",
	"Bhadra",
	"Ashoj",
	"Kartik",
	"Mangsir",
	"Paush",
	"Magh",
	"Falgun",
	"Chaitra",
}

type BSMonth int

const (
	Baisakh BSMonth = 1 + iota
	Jestha
	Ashadh
	Shrawan
	Bhadra
	Ashoj
	Kartik
	Mangsir
	Paush
	Magh
	Falgun
	Chaitra
)

func (m BSMonth) String() string {
	if Baisakh <= m && m <= Chaitra {
		return longBSMonthNames[m-1]
	}
	return "invalid-month"
}

type BSDate struct {
	Year  int
	Month BSMonth
	Day   int
}

func BSDateFromGregorian(gregorian time.Time) BSDate {
	year := gregorian.Year()
	firstDayOfGregorianPlus57 := getDateForFirstDayOfMonth(year+57, Baisakh)

	var nepYear int
	if gregorian.Before(firstDayOfGregorianPlus57) {
		nepYear = year + 56
	} else {
		nepYear = year + 57
	}
	monthStarts := getMonthStarts(nepYear)

	nepMonth := 0 // cannot be zero, this is just a placeholder
	nepDay := 0
	firstDayOfNepYear := monthStarts[0]
	gregorianDayGap := int(math.Floor(gregorian.Sub(firstDayOfNepYear).Hours())) / 24

	monthLengths := getMonthLengthsFromStarts(monthStarts)
	for i, days := range monthLengths {
		if gregorianDayGap < days {
			nepMonth = i + 1
			nepDay = gregorianDayGap + 1
			break
		}

		gregorianDayGap -= days
	}
	assert(Baisakh <= BSMonth(nepMonth) && BSMonth(nepMonth) <= Chaitra, "nepMonth is not in valid range")

	return BSDate{nepYear, BSMonth(nepMonth), nepDay}
}
