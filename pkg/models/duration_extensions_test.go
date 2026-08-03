package models

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestUnmarshalText_DayUnit(t *testing.T) {
	for _, input := range []string{"1d", "1d2h", "2d"} {
		var d Duration
		require.NoError(t, d.UnmarshalText([]byte(input)), "input: %s", input)
	}
}

type rollupConfig struct {
	Rollup struct {
		Time Duration `yaml:"time"`
	} `yaml:"rollup"`
}

func TestUnmarshalYAML_DayUnit(t *testing.T) {
	var cfg rollupConfig
	require.NoError(t, yaml.Unmarshal([]byte("rollup:\n  time: 1d2h\n"), &cfg))
}
