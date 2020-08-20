package version

import (
	"strconv"
	"strings"
)

var (
	Qualifiers = []StringItem{"alpha", "beta", "milestone", "rc", "snapshot", "", "sp"}
	Aliases    = map[string]string{"ga": "", "final": "", "release": "", "cr": "rc"}
)

type ItemType int

const (
	StringType ItemType = iota
	IntType
	ListType
)

type Version struct {
	value string
	Items []Item
}

func NewVersion(v string) (*Version, error) {
	parsedVer, err := parseVersion(v)
	if err != nil {
		return nil, err
	}

	return &Version{
		value: v,
		Items: parsedVer,
	}, nil
}

func (v *Version) String() string {
	return v.value
}

func (v *Version) hasListItem() bool {
	for _, item := range v.Items {
		if item.getType() == ListType {
			return true
		}
	}
	return false
}

func (v1 *Version) Compare(v2 Version) int {
	// Padding version items.
	ListPaddingFlag := (v1.hasListItem() || v2.hasListItem())

	v1Items := v1.Items
	v2Items := v2.Items
	if div := len(v1Items) - len(v2Items); div != 0 {
		if div > 0 {
			for i := 0; i < div; i++ {
				if ListPaddingFlag {
					v2Items = append(v2Items, ListItem{StringItem("")})
				} else {
					v2Items = append(v2Items, IntItem(0))
				}
			}
		}
		if div < 0 {
			for i := div; i < 0; i++ {
				if ListPaddingFlag {
					v1Items = append(v1Items, ListItem{StringItem("")})
				} else {
					v1Items = append(v1Items, IntItem(0))
				}
			}
		}
	}

	for i, item := range v1Items {
		if item.isNull() && v2Items[i].isNull() {
			continue
		}
		result := item.Compare(v2Items[i])
		if result != 0 {
			return result
		}
	}

	return 0
}

func (v1 *Version) Equal(v2 Version) bool {
	return v1.Compare(v2) == 0
}

func (v1 *Version) GreaterThan(v2 Version) bool {
	return v1.Compare(v2) > 0
}

func (v1 *Version) LessThan(v2 Version) bool {
	return v1.Compare(v2) < 0
}

type Item interface {
	Compare(v2 Item) int
	getType() ItemType
	isNull() bool
}

func parseItem(item string) Item {
	i, err := strconv.Atoi(item)
	if err != nil {
		return StringItem(item)
	}
	return IntItem(i)
}

type StringItem string

func (item1 StringItem) Compare(item2 Item) int {
	switch item2.getType() {
	case IntType:
		//  11.alpha < 11.0 < 11.a
		if item2.isNull() {
			if item1.includeWithArray(Qualifiers) < StringItem("").includeWithArray(Qualifiers) {
				return -1
			}
			return 1
		}
		// 11.a < 11.1
		return -1

	case StringType:
		if item1 == item2.(StringItem) {
			return 0
		}

		// "alpha", "beta", "milestone", "rc", "snapshot", "", "sp"
		q1 := item1.includeWithArray(Qualifiers)
		q2 := item2.(StringItem).includeWithArray(Qualifiers)
		if q1 > q2 {
			return 1
		} else if q1 < q2 {
			return -1
		}

		// 11.a < 11.b
		if item1 > item2.(StringItem) {
			return 1
		}
		return -1

	case ListType:
		// 11.a < 11-1
		if len(item2.(ListItem)) == 0 {
			return 1
		}
		return -1
	}
	return 0
}

func (item StringItem) getType() ItemType {
	return StringType
}

func (item StringItem) isNull() bool {
	if item == "" {
		return true
	}
	return false
}

func (item StringItem) includeWithArray(sa []StringItem) int {
	for i, q := range sa {
		if q == item {
			return i
		}
	}
	return len(sa)
}

type IntItem int

func (item1 IntItem) Compare(item2 Item) int {
	switch item2.getType() {
	case IntType:
		if item1 == item2.(IntItem) {
			return 0
		}

		if item1 > item2.(IntItem) {
			return 1
		}
		return -1

	case StringType:
		// 1.alpha < 1
		if item1.isNull() && !item2.(StringItem).isNull() {
			if item2.(StringItem).includeWithArray(Qualifiers) < StringItem("").includeWithArray(Qualifiers) {
				return 1
			}
			return -1
		}
		return 1

	case ListType:
		return 1
	}
	return 0
}

func (item IntItem) getType() ItemType {
	return IntType
}

func (item IntItem) isNull() bool {
	if item == 0 {
		return true
	}
	return false
}

type ListItem []Item

func (item ListItem) getType() ItemType {
	return ListType
}

func (listitem1 ListItem) Compare(item2 Item) int {
	switch item2.getType() {
	case IntType:
		return -1
	case StringType:
		return 1
	case ListType:
		// Padding list items.
		listitem2 := item2.(ListItem)
		if div := len(listitem1) - len(listitem2); div != 0 {
			if div > 0 {
				for i := 0; i < div; i++ {
					listitem2 = append(listitem2, IntItem(1))
				}
			}
			if div < 0 {
				for i := div; i < 0; i++ {
					listitem1 = append(listitem1, IntItem(1))
				}
			}
		}

		for i, item := range listitem1 {
			result := item.Compare(listitem2[i])
			if result != 0 {
				return result
			}
		}
		return 0
	}
	return 0
}

func (items ListItem) isNull() bool {
	if len(items) == 0 {
		return true
	}
	for _, item := range items {
		if !item.isNull() {
			return false
		}
	}
	return true
}

func stringItem(item string) StringItem {
	switch item {
	case "a":
		return StringItem("alpha")
	case "b":
		return StringItem("beta")
	case "m":
		return StringItem("milestone")
	}

	return StringItem(item)
}

// parseVersion is normalize version string.
func parseVersion(v string) ([]Item, error) {
	var stack []Item
	var list ListItem

	isDigit := false
	startIndex := 0
	str := strings.ToLower(v)
	sa := strings.Split(str, "")
	for i, c := range sa {
		if c == "." {
			if i != startIndex {
				s, ok := Aliases[str[startIndex:i]]
				if ok || s != "" {
					list = append(list, parseItem(s))
				} else {
					list = append(list, parseItem(str[startIndex:i]))
				}
			} else {
				list = append(list, IntItem(0))
			}
			startIndex = i + 1
		} else if c == "-" {
			if i != startIndex {
				s, ok := Aliases[str[startIndex:i]]
				if ok || s != "" {
					list = append(list, parseItem(s))
				} else {
					list = append(list, parseItem(str[startIndex:i]))
				}
			} else {
				list = append(list, IntItem(0))
			}
			startIndex = i + 1

			stack = append(stack, list)
			list = ListItem{}

		} else if _, err := strconv.Atoi(c); err == nil {
			if !isDigit && i > startIndex {
				// list = append(list, str[startIndex:i])
				s, ok := Aliases[str[startIndex:i]]
				if ok || s != "" {
					list = append(list, stringItem(s))
				} else {
					list = append(list, stringItem(str[startIndex:i]))
				}
				startIndex = i

				stack = append(stack, list)
				list = ListItem{}
			}

			isDigit = true
		} else {
			if isDigit && i > startIndex {
				// list = append(list, parseItem(str[startIndex:i]))
				list = append(list, parseItem(str[startIndex:i]))
				startIndex = i

				stack = append(stack, list)
				list = ListItem{}
			}
			isDigit = false
		}
	}
	if len(v) > startIndex {
		s, ok := Aliases[str[startIndex:]]
		if ok || s != "" {
			list = append(list, parseItem(s))
		} else {
			list = append(list, parseItem(str[startIndex:]))
		}

		stack = append(stack, list)
	}

	ret := []Item{}
	for _, item := range stack[0].(ListItem) {
		ret = append(ret, item)
	}
	ret = trimNullSuffix(ret)

	return append(ret, stack[1:]...), nil
}

// trimNullSuffix is triming null item.
// ex... (1.0.0 == 1), (1-- == 1), (1.. == 1), (1.0.a.0 == 1.0.a)
func trimNullSuffix(items []Item) []Item {
	ret := items
	for i := len(items) - 1; i >= 0; i-- {
		if items[i].isNull() {
			ret = ret[:i]
		} else {
			break
		}
	}
	return ret
}
