package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

var (
	Qualifiers = []string{"alpha", "beta", "milestone", "rc", "snapshot", "", "sp"}
	Aliases    = map[string]string{"ga": "", "final": "", "release": "", "cr": "rc"}
)

type Version string

func (v Version) String() string {
	lstr := []string{}
	for _, list := range v.parseVersion() {
		if len(list) > 0 {
			lstr = append(lstr, strings.Join(list, "."))
		}
	}
	return strings.Join(lstr, "-")
}

func (v1 Version) Compare(v2 Version) int {
	parsedV1 := v1.parseVersion()
	parsedV2 := v2.parseVersion()

	// Padding
	stackDiv := len(parsedV1) - len(parsedV2)
	if stackDiv > 0 {
		for i := 0; i < stackDiv; i++ {
			parsedV2 = append(parsedV2, []string{})
		}
	} else if stackDiv < 0 {
		for i := 0; i < int(math.Abs(float64(stackDiv))); i++ {
			parsedV1 = append(parsedV1, []string{})
		}
	}
	if len(parsedV1) != len(parsedV2) {
		log.Fatal("padding error")
	}
	for i := range parsedV1 {
		parsedV1[i], parsedV2[i] = paddingArray(parsedV1[i], parsedV2[i], "")
	}

	// Compare
	for i := range parsedV1 {
		if len(parsedV1[i]) != len(parsedV2[i]) {
			log.Fatal("padding item error")
		}
		for j := range parsedV1[i] {
			if result := compare(parsedV1[i][j], parsedV2[i][j]); result != 0 {
				return result
			}
		}
	}

	return 0
}

func paddingArray(a1, a2 []string, paddingStr string) ([]string, []string) {
	stackDiv := len(a1) - len(a2)
	if stackDiv > 0 {
		for i := 0; i < stackDiv; i++ {
			a2 = append(a2, paddingStr)
		}
	} else if stackDiv < 0 {
		for i := 0; i < int(math.Abs(float64(stackDiv))); i++ {
			a1 = append(a1, paddingStr)
		}
	}

	return a1, a2
}

func (v1 Version) Equal(v2 Version) bool {
	return v1.Compare(v2) == 0
}

func (v1 Version) GreaterThan(v2 Version) bool {
	return v1.Compare(v2) > 0
}

func (v1 Version) LessThan(v2 Version) bool {
	return v1.Compare(v2) < 0
}

func (v Version) parseVersion() [][]string {
	var stack [][]string
	var list []string

	isDigit := false
	startIndex := 0
	str := strings.ToLower(string(v))
	sa := strings.Split(str, "")
	for i, c := range sa {
		if c == "." {
			if i != startIndex {
				s, ok := Aliases[str[startIndex:i]]
				if ok || s != "" {
					list = append(list, s)
				} else {
					list = append(list, str[startIndex:i])
				}
			} else {
				list = append(list, "0")
			}
			startIndex = i + 1
		} else if c == "-" {
			if i != startIndex {
				s, ok := Aliases[str[startIndex:i]]
				if ok || s != "" {
					list = append(list, s)
				} else {
					list = append(list, str[startIndex:i])
				}
			} else {
				list = append(list, "0")
			}
			startIndex = i + 1
			stack = append(stack, trimNullValue(list))
			list = []string{}

		} else if _, err := strconv.Atoi(c); err == nil {
			if !isDigit && i > startIndex {
				// list = append(list, str[startIndex:i])
				list = append(list, stringItem(str[startIndex:i]))

				startIndex = i

				stack = append(stack, trimNullValue(list))
				list = []string{}
			}

			isDigit = true
		} else {
			if isDigit && i > startIndex {
				// list = append(list, parseItem(str[startIndex:i]))
				list = append(list, str[startIndex:i])

				startIndex = i

				stack = append(stack, trimNullValue(list))
				list = []string{}
			}
			isDigit = false
		}
	}
	if len(v) > startIndex {
		s, ok := Aliases[str[startIndex:]]
		if ok || s != "" {
			list = append(list, s)
		} else {
			list = append(list, str[startIndex:])
		}

		stack = append(stack, trimNullValue(list))
	}

	return stack
}

func stringItem(item string) string {
	switch item {
	case "a":
		return "alpha"
	case "b":
		return "beta"
	case "m":
		return "milestone"
	}

	return item
}

func compare(c1, c2 string) int {
	//  Aliases transfer
	_, ok := Aliases[c1]
	if ok {
		c1 = Aliases[c1]
	}
	_, ok = Aliases[c2]
	if ok {
		c2 = Aliases[c2]
	}

	// Qualifiers compare
	// Qualifiers = []string{"alpha", "beta", "milestone", "rc", "snapshot", "", "sp"}
	// "alpha" < "beta" < "milestone" < "rc" < "snapshot" < "" < "sp" < "[ASCII]" < [Integer]}
	// "1.foo" < "1-foo" < "1-1" < "1.1"
	q1 := includeWithArray(Qualifiers, c1)
	q2 := includeWithArray(Qualifiers, c2)
	if q1 > q2 {
		return 1
	}
	if q1 < q2 {
		return -1
	}

	iFlag1 := false
	i1, err := strconv.Atoi(c1)
	if err == nil {
		iFlag1 = true
	}

	iFlag2 := false
	i2, err := strconv.Atoi(c2)
	if err == nil {
		iFlag2 = true
	}

	// Compare
	if iFlag1 && !iFlag2 {
		return 1
	} else if !iFlag1 && iFlag2 {
		return -1
	} else if iFlag1 && iFlag2 {
		if i1 > i2 {
			return 1
		}
		if i1 < i2 {
			return -1
		}
	} else {
		if c1 > c2 {
			return 1
		}
		if c1 < c2 {
			return -1
		}
	}

	return 0
}

// if not include "s" return "sa" length
func includeWithArray(sa []string, s string) int {
	for i, q := range sa {
		if q == s {
			return i
		}
	}
	return len(sa)
}

// Null value is 0 or ""
func trimNullValue(list []string) []string {
	for i := len(list) - 1; i > -1; i-- {
		if list[i] == "0" || list[i] == "" {
			list = list[:i]
		} else {
			break
		}
	}
	return list
}

func main() {
	args := os.Args
	if len(args) != 3 {
		log.Fatal("args required 2 versions...")
	}
	v1 := Version(args[1])
	v2 := Version(args[2])

	if v1.GreaterThan(v2) {
		fmt.Printf("%s greater than %s\n", v1, v2)
	}
	if v1.LessThan(v2) {
		fmt.Printf("%s less than %s\n", v1, v2)
	}
	if v1.Equal(v2) {
		fmt.Printf("%s equal %s\n", v1, v2)
	}
}
