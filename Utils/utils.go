package Utils

import (
	"strings"
)

// A list of incorrect -> correct name mappings, use this as last resort to
// fix incorrect/inconsistent names
var NameAdjustments = map[string]string{}

func FixName(c string) string {
	vals := strings.Split(strings.TrimSpace(c), " ")
	first := vals[0]
	lastidx := len(vals) - 1
	for (strings.HasPrefix(vals[lastidx], "(") && strings.HasSuffix(vals[lastidx], ")")) || vals[lastidx] == "Jr" || vals[lastidx] == "Jr." || vals[lastidx] == "Sr" || vals[lastidx] == "Sr." || vals[lastidx] == "I" || vals[lastidx] == "II" || vals[lastidx] == "III" || vals[lastidx] == "" {
		lastidx--
	}
	last := vals[lastidx]
	last = strings.TrimSuffix(last, ",")
	name := strings.Title(strings.ToLower(first + " " + last))
	if lastidx == 0 {
		name = strings.Title(strings.ToLower(first))
	}
	if n, ok := NameAdjustments[name]; ok {
		name = n
	}
	return name
}
