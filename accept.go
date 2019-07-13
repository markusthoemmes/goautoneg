package goautoneg

import (
	"sort"
	"strconv"
	"strings"
)

// Accept is a structure to represent a clause in an HTTP Accept Header.
type Accept struct {
	Type, SubType string
	Q             float64
	Params        map[string]string
}

// ParseAccept parses the given string as an Accept header as defined in
// https://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.1.
// Some rules are only loosely applied and might not be as strict as defined in the RFC.
func ParseAccept(header string) []Accept {
	parts := strings.Split(header, ",")
	clauses := []Accept{}

	for _, part := range parts {
		part := trim(part)

		// media-range is defined as
		// media-range = ( "*/*" | ( type "/" "*" ) | ( type "/" subtype )) *( ";" parameter )
		mediaRangeParts := strings.Split(part, ";")

		accept := Accept{
			Q:      1.0, // "[...] The default value is q=1"
			Params: make(map[string]string, len(mediaRangeParts)-1),
		}

		// The type part of the media-range is defined as
		// "*/*" | ( type "/" "*" ) | ( type "/" subtype )
		types := strings.Split(mediaRangeParts[0], "/")

		switch {
		// This case is not defined in the spec keep it to mimic the original code.
		case len(types) == 1 && types[0] == "*":
			accept.Type = "*"
			accept.SubType = "*"
		case len(types) == 2:
			accept.Type = trim(types[0])
			accept.SubType = trim(types[1])
		default:
			continue
		}

		// The parameter part of the media-range is defined as
		// "q" "=" qvalue *( ";" token [ "=" ( token | quoted-string ) )
		for _, param := range mediaRangeParts[1:] {
			paramParts := strings.SplitN(param, "=", 2)
			if len(paramParts) != 2 {
				// Ignore parameters with no delimiter.
				continue
			}

			key := trim(paramParts[0])
			if key == "q" {
				// A parsing failure will set Q to 0.
				accept.Q, _ = strconv.ParseFloat(paramParts[1], 64)
			} else {
				accept.Params[key] = trim(paramParts[1])
			}
		}

		clauses = append(clauses, accept)
	}

	sort.SliceStable(clauses, func(i, j int) bool {
		return clauses[i].Q > clauses[j].Q
	})

	return clauses
}

func trim(s string) string {
	return strings.Trim(s, " ")
}
