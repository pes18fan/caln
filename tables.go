package caln

// For brevity, I shall refer to Bikram Sambat as BS throughout.

import (
	"caln/julianday"
	"fmt"
	"math"
	"time"
)

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

	for i := range 12 {
		target := float64(i) * (math.Pi / 6) // i times 30 degrees
		sankranti, jd := getDateForLongitude(lo, hi, target)

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
	sankranti, _ := getDateForLongitude(lo, hi, 0)
	monthStarts[12] = sankranti
	return monthStarts
}

// Get the BS date and Julian day corresponding to a certain sidereal longitude
// via bisection, given two lo and hi Julian day values known to be just before
// and just after the moment of that longitude.
func getDateForLongitude(lo, hi, lon float64) (time.Time, float64) {
	f := func(x float64) float64 {
		// We need to wrap the sidereal longitude function with another, as
		// it does not return negative values.
		// We therefore interpret longitudes greater than pi as negative.
		l := sunSiderealLongitude(x) - lon
		if l > math.Pi {
			l = l - 2*math.Pi
		}
		return l
	}
	jd := bisect(lo, hi, f)

	ktm, err := time.LoadLocation("Asia/Kathmandu")
	if err != nil {
		panic(fmt.Sprintf("failed to load Asia/Kathmandu location: %v", err))
	}
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
	return time.Date(sankranti.Year(), sankranti.Month(),
		sankranti.Day(), 0, 0, 0, 0, sankranti.Location()), jd
}

func getMonthLengths(bsYear int) [12]int {
	monthStarts := getMonthStarts(bsYear)
	return getMonthLengthsFromStarts(monthStarts)
}

func getMonthLengthsFromStarts(monthStarts [13]time.Time) [12]int {
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
