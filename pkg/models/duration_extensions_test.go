package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

// TestUnmarshalText_DayUnit proves the bug: the Groundcover backend normalises any
// rollup duration ≥24h to the "d" unit (e.g. 26h → "1d2h"). time.ParseDuration
// rejects "d", so YAML unmarshal of monitor responses always fails for those durations.
// With strfmt.ParseDuration (already a direct dependency) the tests pass.
func TestUnmarshalText_DayUnit(t *testing.T) {
	cases := []struct {
		input string
		want  time.Duration
	}{
		{"1d", 24 * time.Hour},
		{"1d2h", 26 * time.Hour},
		{"2d", 48 * time.Hour},
	}
	for _, tc := range cases {
		var d Duration
		require.NoError(t, d.UnmarshalText([]byte(tc.input)), "input: %s", tc.input)
		assert.Equal(t, Duration(tc.want), d, "input: %s", tc.input)
	}
}

type rollupConfig struct {
	Rollup struct {
		Time Duration `yaml:"time"`
	} `yaml:"rollup"`
}

// TestUnmarshalYAML_DayUnit reproduces the exact failure seen in terraform-provider-groundcover:
// yaml.Unmarshal of a monitor GET response fails when rollup.time contains "1d2h".
func TestUnmarshalYAML_DayUnit(t *testing.T) {
	input := "rollup:\n  time: 1d2h\n"
	var cfg rollupConfig
	require.NoError(t, yaml.Unmarshal([]byte(input), &cfg))
	assert.Equal(t, Duration(26*time.Hour), cfg.Rollup.Time)
}
