package caln

import (
	_ "embed"
	"strconv"
	"strings"
)

//go:embed VSOP87D.ear.clean
var vsop87d string

type vsopCoefficients [3]float64

// L0 terms to L5 terms, six term lists
var vsopTable [6][]vsopCoefficients

func initVsopTable() {
	vsopTable = calculateVsopTable()
}

func calculateVsopTable() [6][]vsopCoefficients {
	s := strings.TrimSpace(vsop87d)
	lines := strings.Split(s, "\n")

	var table [6][]vsopCoefficients
	for i := range table {
		table[i] = make([]vsopCoefficients, 0)
	}

	for _, line := range lines {
		fields := strings.Split(line, " ")

		coordinateCode := fields[0]
		assert(len(coordinateCode) == 4, "coordinate code is not four bytes long")

		if coordinateCode[2] != '1' {
			continue
		}
		powerChr := coordinateCode[3]
		assert(int(powerChr) > 47 && int(powerChr) < 58, "power is not a valid digit")
		power := byte(int(powerChr) - 48)
		assert(power < 6, "power is not within [0, 5]")

		A, err1 := strconv.ParseFloat(fields[len(fields)-3], 64)
		assert(err1 == nil, "coefficient A is not a valid float")
		B, err2 := strconv.ParseFloat(fields[len(fields)-2], 64)
		assert(err2 == nil, "coefficient B is not a valid float")
		C, err3 := strconv.ParseFloat(fields[len(fields)-1], 64)
		assert(err3 == nil, "coefficient C is not a valid float")

		coeffs := vsopCoefficients{A, B, C}
		table[power] = append(table[power], coeffs)
	}

	return table
}
