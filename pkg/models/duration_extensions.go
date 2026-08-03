// Custom extensions for generated models
// This file contains custom methods for generated types that need additional functionality.
// It gets copied to sdk/pkg/models during the swagger generation process.

package models

import (
	"fmt"
	"regexp"
	"time"

	"github.com/go-openapi/strfmt"
)

// Custom methods for Duration type
// These methods provide YAML and text unmarshalling capabilities

// String returns the string representation of the duration
func (d Duration) String() string {
	return time.Duration(d).String()
}

// validDuration matches only complete duration strings with integer magnitudes and
// known units (no leading sign, no decimals, no trailing garbage).
var validDuration = regexp.MustCompile(`^(\d+(ns|us|µs|ms|s|m|h|d|w))+$`)

// UnmarshalText implements the text unmarshaller interface
func (d *Duration) UnmarshalText(text []byte) error {
	s := string(text)
	if !validDuration.MatchString(s) {
		return fmt.Errorf("invalid duration %q", s)
	}
	duration, err := strfmt.ParseDuration(s)
	if err != nil {
		return err
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
