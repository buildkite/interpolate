package interpolate

import (
	"strings"

	"github.com/drone/envsubst/path"
)

// Imported from drone.io:
// github.com/drone/envsubst/funcs.go

func trimShortestPrefix(s string, args ...string) string {
	if len(args) != 0 {
		s = trimShortest(s, args[0])
	}
	return s
}

func trimShortestSuffix(s string, args ...string) string {
	if len(args) != 0 {
		r := reverse(s)
		rarg := reverse(args[0])
		s = reverse(trimShortest(r, rarg))
	}
	return s
}

func trimLongestPrefix(s string, args ...string) string {
	if len(args) != 0 {
		s = trimLongest(s, args[0])
	}
	return s
}

func trimLongestSuffix(s string, args ...string) string {
	if len(args) != 0 {
		r := reverse(s)
		rarg := reverse(args[0])
		s = reverse(trimLongest(r, rarg))
	}
	return s
}

func trimShortest(s, arg string) string {
	var shortestMatch string
	for i := 0; i < len(s); i++ {
		match, err := path.Match(arg, s[0:len(s)-i])

		if err != nil {
			return s
		}

		if match {
			shortestMatch = s[0 : len(s)-i]
		}
	}

	if shortestMatch != "" {
		return strings.TrimPrefix(s, shortestMatch)
	}

	return s
}

func trimLongest(s, arg string) string {
	for i := 0; i < len(s); i++ {
		match, err := path.Match(arg, s[0:len(s)-i])

		if err != nil {
			return s
		}

		if match {
			return strings.TrimPrefix(s, s[0:len(s)-i])
		}
	}

	return s
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
