package models

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestDurationUnmarshalTextAcceptsPrometheusAndStdlibDurations(t *testing.T) {
	tests := []struct {
		name string
		text string
		want time.Duration
	}{
		{name: "prometheus day", text: "1d", want: 24 * time.Hour},
		{name: "prometheus day and hours", text: "1d2h", want: 26 * time.Hour},
		{name: "prometheus week", text: "1w", want: 7 * 24 * time.Hour},
		{name: "prometheus year", text: "1y", want: 365 * 24 * time.Hour},
		{name: "prometheus 400 days", text: "400d", want: 400 * 24 * time.Hour},
		{name: "prometheus seconds and milliseconds", text: "1s500ms", want: 1500 * time.Millisecond},
		{name: "prometheus minutes seconds and milliseconds", text: "1m30s500ms", want: time.Minute + 30*time.Second + 500*time.Millisecond},
		{name: "prometheus compound", text: "1y2w3d4h5m6s7ms", want: (365+14+3)*24*time.Hour + 4*time.Hour + 5*time.Minute + 6*time.Second + 7*time.Millisecond},
		{name: "prometheus negative day", text: "-1d", want: -24 * time.Hour},
		{name: "prometheus negative day and hours", text: "-1d2h", want: -26 * time.Hour},
		{name: "stdlib decimal hours", text: "1.5h", want: 90 * time.Minute},
		{name: "stdlib leading decimal hours", text: "0.5h", want: 30 * time.Minute},
		{name: "stdlib negative hours", text: "-1h", want: -time.Hour},
		{name: "stdlib nanosecond", text: "1ns", want: time.Nanosecond},
		{name: "stdlib ascii microsecond", text: "1us", want: time.Microsecond},
		{name: "stdlib micro sign microsecond", text: "1µs", want: time.Microsecond},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Duration
			require.NoError(t, d.UnmarshalText([]byte(tt.text)))
			require.Equal(t, tt.want, time.Duration(d))
		})
	}
}

func TestDurationUnmarshalTextRejectsInvalidDurations(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{name: "decimal day", text: "1.5d"},
		{name: "leading decimal day", text: ".5d"},
		{name: "negative decimal day", text: "-1.5d"},
		{name: "negative leading decimal day", text: "-.5d"},
		{name: "unsupported minute alias", text: "1min"},
		{name: "malformed day suffix", text: "1dayzzz"},
		{name: "trailing plus", text: "1d+"},
		{name: "duration overflow", text: "1000000w"},
		{name: "non-leading sign", text: "1d-2h"},
		{name: "day hour rewrite int overflow", text: strconv.Itoa(maxPromHours/promHours["d"]+1) + "d"},
		{name: "week hour rewrite int overflow", text: strconv.Itoa(maxPromHours/promHours["w"]+1) + "w"},
		{name: "year hour rewrite int overflow", text: strconv.Itoa(maxPromHours/promHours["y"]+1) + "y"},
		{name: "negative day hour rewrite int overflow", text: strconv.Itoa(minPromHours/promHours["d"]-1) + "d"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Duration
			require.Error(t, d.UnmarshalText([]byte(tt.text)))
		})
	}
}

func TestDurationUnmarshalYAMLAcceptsPrometheusDuration(t *testing.T) {
	var target struct {
		Rollup struct {
			Time Duration `yaml:"time"`
		} `yaml:"rollup"`
	}

	require.NoError(t, yaml.Unmarshal([]byte("rollup:\n  time: 1d2h\n"), &target))
	require.Equal(t, 26*time.Hour, time.Duration(target.Rollup.Time))
}

func TestDurationMarshalYAMLRoundTrip(t *testing.T) {
	original := Duration(90500 * time.Millisecond)
	marshaled, err := original.MarshalYAML()
	require.NoError(t, err)

	var decoded Duration
	require.NoError(t, decoded.UnmarshalText([]byte(marshaled.(string))))
	require.Equal(t, time.Duration(original), time.Duration(decoded))
}
