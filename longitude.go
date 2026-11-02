package caln

import (
	"math"
)

// approximation for the ayanamsa for a given julian day
// Unit is radians
func ayanamsa(jd float64) float64 {
	ayanamsaAtJ2000 := 23.853
	rateOfPrecession := 0.01397
	return degreesToRadians(ayanamsaAtJ2000 + rateOfPrecession*((jd-2451545.0)/365.25))
}

// Returns value in radians
func sunTropicalLongitude(jd float64) float64 {
	// calculate julian millenia since J2000.0
	t := (jd - 2451545) / 365250

	// calculate Earth's heliocentric longitude
	var earthHeliocentricLongitude float64
	for i, series := range vsopTable {
		var seriesSum float64
		for _, coeffs := range series {
			A, B, C := coeffs[0], coeffs[1], coeffs[2]
			seriesSum += A * math.Cos(B+C*t)
		}
		earthHeliocentricLongitude += seriesSum * math.Pow(t, float64(i))
	}

	// calculate the Sun's geocentric longitude, and normalize it to be within
	// [0, 2pi]
	sunGeocentricLongitude := math.Mod(earthHeliocentricLongitude+math.Pi, 2*math.Pi)

	// apply corrections to get the apparent tropical longitude
	// NOTE: the corrections are just approximations
	T := (jd - 2451545) / 36525 // julian centuries since J2000.0
	return sunGeocentricLongitude -
		9.930923e-05 - 8.342674e-05*math.Sin(2.1823597-33.7570414*T)
}

// Sidereal longitude is calculated by first calculating the tropical one then
// subtracting the Chitrapaksha ayanamsa
// The returned value is in radians
func sunSiderealLongitude(jd float64) float64 {
	tropical := sunTropicalLongitude(jd)
	sidereal := math.Mod(tropical-ayanamsa(jd), 2*math.Pi)
	if sidereal < 0 {
		sidereal += 2 * math.Pi
	}
	return sidereal
}
