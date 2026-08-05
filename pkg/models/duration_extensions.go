// Custom extensions for generated models
// This file contains custom methods for generated types that need additional functionality.
// It gets copied to sdk/pkg/models during the swagger generation process.

package models

import (
	"regexp"
	"strconv"
	"time"
)

// Custom methods for Duration type
// These methods provide YAML and text unmarshalling capabilities

// String returns the string representation of the duration
func (d Duration) String() string {
	return time.Duration(d).String()
}

// Some server-side durations, notably monitor rollup.time, are prometheus
// model.Duration values, so they arrive in prometheus format ("1d2h", "1y").
// The SDK generator emits one shared Duration type, so this parser intentionally
// accepts both stdlib duration strings and prometheus day/week/year units.
// Rewrite the units time.ParseDuration doesn't know into hours; the stdlib does
// all final parsing and validation.
var promHours = map[string]int{"y": 365 * 24, "w": 7 * 24, "d": 24}

var maxPromHours = int(^uint(0) >> 1)
var minPromHours = -maxPromHours - 1

// The optional decimal groups are deliberate: they make "1.5d" and ".5d" match
// whole so Atoi rejects them, instead of the regex matching just "5d" and
// silently yielding 1.12h or .120h.
var promUnitRE = regexp.MustCompile(`(-?(?:\d+(?:\.\d+)?|\.\d+))([ydw])`)

// UnmarshalText implements the text unmarshaller interface. It accepts Go stdlib
// duration strings and the prometheus day/week/year units emitted by server-side
// model.Duration fields.
func (d *Duration) UnmarshalText(text []byte) error {
	duration, err := time.ParseDuration(string(text))
	if err != nil {
		s := promUnitRE.ReplaceAllStringFunc(string(text), func(m string) string {
			n, aerr := strconv.Atoi(m[:len(m)-1])
			h := promHours[m[len(m)-1:]]
			if aerr != nil || n > maxPromHours/h || n < minPromHours/h {
				return m
			}
			return strconv.Itoa(n*h) + "h"
		})
		if duration, err = time.ParseDuration(s); err != nil {
			return err
		}
	}
	*d = Duration(duration)
	return nil
}

// UnmarshalYAML implements the YAML unmarshaller interface
func (d *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	return d.UnmarshalText([]byte(s))
}

// MarshalYAML implements the YAML marshaller interface
func (d Duration) MarshalYAML() (interface{}, error) {
	return d.String(), nil
}
