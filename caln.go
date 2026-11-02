package caln

// For brevity, I shall refer to Bikram Sambat as BS throughout.

import (
	"caln/julianday"
	"fmt"
	"math"
	"time"
)

func init() {
	initVsopTable()
}

// Get the time.Times that correspond to the first day of each BS month of
// bsYear, plus the first day of Baisakh the year after.
func getMonthStarts(bsYear int) [13]time.Time {
	// We need the Julian day for the beginning of each BS month.
	// What is the Julian day? It is somewhat similar in concept to the Unix
	// timestamp, but unlike the latter which counts the number of seconds since
	// January 1 1970, it counts the number of DAYS passed since the Julian
	// period (4713 BC).
	// But how do we get that for the first day of a BS month, say Baisakh
	// 1, when we don't even have a BS date yet? This might sound like a
	// chicken and egg situation but it is possible to get this value through
	// approximations.
	// Baisakh 1, which corresponds to Mesha Sankranti, is the specific day
	// when the Sun's sidereal longitude is zero.
	// We have a function that can calculate the sidereal longitude for a
	// certain Julian day, but not a function that can do it the other way
	// around. Thus, we need to use the numerical method for bisection.
	// We select two Julian days that occur just before and just after Mesha
	// Sankranti i.e. give us negative and positive values for the sidereal
	// longitude function respectively, and then repeatedly take the midpoint
	// from there till we are satisfied with the accuracy.
	// And we can then generalize this for each month. Jestha for example starts
	// when the Sun's sidereal longitude is thirty degrees or pi/6 radians, we
	// can bisect around that the same way. For Ashadh it is sixty degrees,
	// and so on.

	// Beginning of each month of the current year, plus the beginning of Baisakh
	// of next year.
	var monthStarts [13]time.Time

	// These are the values used for bisection of Baisakh 1. We can reuse the
	// lo and hi values updating them as they go to avoid recomputing
	// everything, which is why I'm defining them outside the loop.
	lo, hi :=
		// Baisakh 1 is generally around April 13-15, so picking April 10 as the
		// low bound and April 18 as the high bound is good enough.
		float64(julianday.Number(bsYear-57, time.April, 10)),
		float64(julianday.Number(bsYear-57, time.April, 18))

	// We need this funny little function-creating-function because we need a
	// function to pass into the bisector, but that function also needs to
	// calculate different results based on the target we want to bisect at.
	fMaker := func(target float64) func(float64) float64 {
		return func(x float64) float64 {
			// We need to wrap the sidereal longitude function with another, as
			// it does not return negative values.
			// We therefore interpret longitudes greater than pi as negative.
			l := sunSiderealLongitude(x) - target
			if l > math.Pi {
				l = l - 2*math.Pi
			}
			return l
		}
	}

	ktm, err := time.LoadLocation("Asia/Kathmandu")
	if err != nil {
		panic(fmt.Sprintf("failed to load Asia/Kathmandu location: %v", err))
	}
	for i := range 12 {
		target := float64(i) * (math.Pi / 6) // i times 30 degrees
		f := fMaker(target)
		jd := bisect(lo, hi, f)

		// We have the Julian day for the start of this month! We can now
		// convert it to a localized time.Time (which is done by ToTime through
		// conversion into the Unix timestamp).
		sankranti := julianday.ToTime(jd).In(ktm)

		// If the Sankranti falls after sunset, Nepal assigns the new month
		// to the next day. So, we need to account for that.
		sunset := sunsetTimeForKathmandu(sankranti)
		if sankranti.After(sunset) {
			sankranti = sankranti.AddDate(0, 0, 1)
		}

		// Remove the hour, minute, second etc components, we don't need that
		// anymore.
		sankranti = time.Date(sankranti.Year(), sankranti.Month(),
			sankranti.Day(), 0, 0, 0, 0, sankranti.Location())

		// Update lo and hi. lo can be the Julian day we just calculated, and
		// hi can be, about 40 days past that (as that is further than
		// the start of next month); forming a valid bisection interval for
		// the next calculation.
		lo = jd
		hi = jd + 40

		monthStarts[i] = sankranti
	}

	// To calculate Chaitra's length we need to repeat that above loop
	// calculation just one more time for the target 0 (next Baisakh), with the
	// accumulated lo and hi.
	jd := bisect(lo, hi, fMaker(0))
	sankranti := julianday.ToTime(jd).In(ktm)
	sunset := sunsetTimeForKathmandu(sankranti)
	if sankranti.After(sunset) {
		sankranti = sankranti.AddDate(0, 0, 1)
	}
	sankranti = time.Date(sankranti.Year(), sankranti.Month(),
		sankranti.Day(), 0, 0, 0, 0, sankranti.Location())
	monthStarts[12] = sankranti
	return monthStarts
}

func getMonthLengths(bsYear int) [12]int {
	monthStarts := getMonthStarts(bsYear)
	return getMonthLengthsFromStarts(bsYear, monthStarts)
}

func getMonthLengthsFromStarts(bsYear int, monthStarts [13]time.Time) [12]int {
	var monthLengths [12]int
	for i := range 12 {
		prev, next := monthStarts[i], monthStarts[i+1]
		prevYear, nextYear := prev.Year(), next.Year()
		if prevYear == nextYear {
			monthLengths[i] = next.YearDay() - prev.YearDay()
		} else {
			// In this case `next` is in the next year
			yearDays := 365
			if isLeapYear(prevYear) {
				yearDays = 366
			}
			monthLengths[i] = yearDays + next.YearDay() - prev.YearDay()
		}

	}

	return monthLengths
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
	monthStarts := getMonthStarts(year + 57)

	var nepYear int
	if gregorian.Before(monthStarts[0]) {
		monthStarts = getMonthStarts(year + 56)
		nepYear = year + 56
	} else {
		nepYear = year + 57
	}

	nepMonth := 0
	nepDay := 0
	firstDayOfNepYear := monthStarts[0]
	gregorianDayGap := int(math.Floor(gregorian.Sub(firstDayOfNepYear).Hours())) / 24

	monthLengths := getMonthLengthsFromStarts(nepYear, monthStarts)
	for i, days := range monthLengths {
		if gregorianDayGap < days {
			nepMonth = i + 1
			nepDay = gregorianDayGap + 1
			break
		}

		gregorianDayGap -= days
	}

	return BSDate{nepYear, BSMonth(nepMonth), nepDay}
}
