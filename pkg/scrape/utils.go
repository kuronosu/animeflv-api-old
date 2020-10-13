package scrape

import (
	"bytes"
	"math"
	"strconv"
)

// EmailProtected represents [email protected] string
var EmailProtected = parseEmailProtectedChars()

func parseEmailProtectedChars() string {
	emailProtectedChars := []int{91, 101, 109, 97, 105, 108, 160, 112, 114, 111, 116, 101, 99, 116, 101, 100, 93}
	ep := ""
	for _, asciiNum := range emailProtectedChars {
		ep += string(asciiNum)
	}
	return ep
}

// DecodeEmail decode email from cloudflare
func DecodeEmail(a string) string {
	var e bytes.Buffer
	r, _ := strconv.ParseInt(a[0:2], 16, 0)
	for n := 4; n < len(a)+2; n += 2 {
		i, _ := strconv.ParseInt(a[n-2:n], 16, 0)
		e.WriteString(string(i ^ r))
	}
	return e.String()
}

// MakeRange asd
func MakeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func appendCategory(a []string, b []string) []string {
	check := make(map[string]int)
	d := append(a, b...)
	res := make([]string, 0)
	for _, val := range d {
		check[val] = 1
	}

	for letter := range check {
		res = append(res, letter)
	}

	return res
}

// FloatToString convert float64 to string
func FloatToString(num float64) string {
	decimals := 1
	if math.Trunc(float64(num)) == float64(num) {
		decimals = 0
	}
	return strconv.FormatFloat(float64(num), 'f', decimals, 64)
}

// ValidLang verify if string is a valid lang
func ValidLang(lang string) bool {
	switch lang {
	case SUB, LAT, ESP:
		return true
	}
	return false
}

// AvailableServers programed scrape servers
var AvailableServers = []string{"gocdn"}

// ValidServer verify if string is a valid server
func ValidServer(server string) bool {
	for _, s := range AvailableServers {
		if s == server {
			return true
		}
	}
	return false
}
