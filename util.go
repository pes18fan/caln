package caln

import (
	"math"
	"time"
)

// NOAA algorithm
// All kinda mumbo jumbo for me
func sunsetTimeForKathmandu(t time.Time) time.Time {
	lat := 27.7172
	lon := 85.3240

	yearDays := 365
	if isLeapYear(t.Year()) {
		yearDays = 366
	}
	y := 2 * math.Pi / float64(yearDays) * (float64(t.YearDay()) - 1) // fractional year
	eqTime := 229.18 * (0.000075 +
		0.001868*math.Cos(y) -
		0.032077*math.Sin(y) -
		0.014615*math.Cos(2*y) -
		0.040849*math.Sin(2*y))
	declination := 0.006918 -
		0.399912*math.Cos(y) + 0.070257*math.Sin(y) - 0.006758*math.Cos(2*y) +
		0.000907*math.Sin(2*y) - 0.002697*math.Cos(3*y) + 0.00148*math.Sin(3*y)
	hourAngle := math.Acos(
		math.Cos(1.585335)/
			(math.Cos(degreesToRadians(lat))*math.Cos(declination)) -
			(math.Tan(degreesToRadians(lat)) * math.Tan(declination)))
	sunsetUTCMinutes := 720 - 4*(lon-radiansToDegrees(hourAngle)) - eqTime
	duration := time.Duration(sunsetUTCMinutes * float64(time.Minute))

	base := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	return base.Add(duration)
}

func degreesToRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}

func radiansToDegrees(rad float64) float64 {
	return rad * (180 / math.Pi)
}

func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

const bisection_iterations = 50

func bisect(x, y float64, f func(float64) float64) float64 {
	if f(x) > 0 || f(y) < 0 {
		panic("starting values for bisection are not valid")
	}

	var z float64
	for range bisection_iterations {
		z = (x + y) / 2
		if f(z) < 0 {
			x = z
		} else {
			y = z
		}
	}
	return z
}

func assert(cond bool, message string) {
	if !cond {
		panic(message)
	}
}
