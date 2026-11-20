package orbitalwiki

import (
	"fmt"
	"strings"
)

// OrbitClass is the API's orbit_class enum. The zero value means "not set".
type OrbitClass int

const (
	OrbitClassUnset OrbitClass = iota
	OrbitClassLEO
	OrbitClassMEO
	OrbitClassGEO
	OrbitClassHEO
	OrbitClassSSO
	OrbitClassVLEO
	OrbitClassUnknown
)

// orbitClassNames is indexed by OrbitClass.
var orbitClassNames = []string{"", "LEO", "MEO", "GEO", "HEO", "SSO", "VLEO", "UNKNOWN"}

// String returns the API value ("LEO", …), or "" for OrbitClassUnset.
func (o OrbitClass) String() string { return enumString(orbitClassNames, int(o)) }

// ParseOrbitClass parses an API value case-insensitively, ignoring surrounding whitespace.
// "" parses to OrbitClassUnset with no error. Any other unrecognised value returns an
// error that lists the allowed values.
func ParseOrbitClass(s string) (OrbitClass, error) {
	i, err := parseEnum("orbit class", orbitClassNames, s)
	return OrbitClass(i), err
}

// OrbitClassValues returns every API value, excluding Unset, in declaration order.
func OrbitClassValues() []string { return enumValues(orbitClassNames) }

// Status is the API's operational status enum. The zero value means "not set".
type Status int

const (
	StatusUnset Status = iota
	StatusOperational
	StatusNonOperational
	StatusPartiallyOperational
	StatusStandby
	StatusSpare
	StatusExtendedMission
	StatusDecayed
	StatusUnknown
)

// statusNames is indexed by Status.
var statusNames = []string{
	"", "OPERATIONAL", "NONOPERATIONAL", "PARTIALLY_OPERATIONAL", "STANDBY",
	"SPARE", "EXTENDED_MISSION", "DECAYED", "UNKNOWN",
}

// String returns the API value ("OPERATIONAL", …), or "" for StatusUnset.
func (s Status) String() string { return enumString(statusNames, int(s)) }

// ParseStatus parses an API value case-insensitively, ignoring surrounding whitespace.
// "" parses to StatusUnset with no error. Any other unrecognised value returns an
// error that lists the allowed values.
func ParseStatus(s string) (Status, error) {
	i, err := parseEnum("status", statusNames, s)
	return Status(i), err
}

// StatusValues returns every API value, excluding Unset, in declaration order.
func StatusValues() []string { return enumValues(statusNames) }

func enumString(names []string, i int) string {
	if i < 0 || i >= len(names) {
		return ""
	}
	return names[i]
}

func parseEnum(kind string, names []string, s string) (int, error) {
	v := strings.ToUpper(strings.TrimSpace(s))
	for i, name := range names {
		if v == name {
			return i, nil
		}
	}
	return 0, fmt.Errorf("invalid %s %q: must be one of %s", kind, s, strings.Join(names[1:], ", "))
}

func enumValues(names []string) []string {
	return append([]string(nil), names[1:]...)
}
