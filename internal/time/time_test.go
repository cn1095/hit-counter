package time

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimestampByMaxTime(t *testing.T) {
	assert.Equal(t, int64((1<<63)-1), TimestampByMaxTime())
}

func TestStringToTime(t *testing.T) {
	now := time.Unix(time.Now().Unix(), 0).UTC()
	tests := []struct {
		name   string
		input  string
		output time.Time
	}{
		{
			name:   "empty",
			output: time.Time{},
		},
		{
			name: "exist",
			input: fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", now.Year(), now.Month(), now.Day(),
				now.Hour(), now.Minute(), now.Second()),
			output: now,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.output, StringToTime(tc.input))
		})
	}
}

func TestTimeToString(t *testing.T) {
	now := time.Unix(time.Now().Unix(), 0).UTC()
	tests := []struct {
		name   string
		input  time.Time
		output string
	}{
		{
			name: "exist",
			output: fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", now.Year(), now.Month(), now.Day(),
				now.Hour(), now.Minute(), now.Second()),
			input: now,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.output, TimeToString(tc.input))
		})
	}
}

func TestYearlyStringToTime(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output time.Time
	}{
		{
			name:   "empty",
			output: time.Time{},
		},
		{
			name:   "exist",
			input:  fmt.Sprintf("%d", time.Now().Year()),
			output: time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.output, YearlyStringToTime(tc.input))
		})
	}
}

func TestTimeToYearlyStringFormat(t *testing.T) {
	tests := []struct {
		name   string
		input  time.Time
		output string
	}{
		{
			name:   "exist",
			input:  time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.UTC),
			output: fmt.Sprintf("%d", time.Now().Year()),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.output, TimeToYearlyStringFormat(tc.input))
		})
	}
}

func TestMonthlyStringToTime(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output time.Time
	}{
		{
			name:   "empty",
			output: time.Time{},
		},
		{
			name:   "exist",
			input:  fmt.Sprintf("%d%02d", time.Now().Year(), time.Now().Month()),
			output: time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.output, MonthlyStringToTime(tc.input))
		})
	}
}

func TestTimeToMonthlyStringFormat(t *testing.T) {
	tests := []struct {
		name   string
		input  time.Time
		output string
	}{
		{
			name:   "exist",
			input:  time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC),
			output: fmt.Sprintf("%d%02d", time.Now().Year(), time.Now().Month()),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.output, TimeToMonthlyStringFormat(tc.input))
		})
	}
}

func TestDailyStringToTime(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output time.Time
	}{
		{
			name:   "empty",
			output: time.Time{},
		},
		{
			name:   "exist",
			input:  fmt.Sprintf("%d%02d%02d", time.Now().Year(), time.Now().Month(), time.Now().Day()),
			output: time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.output, DailyStringToTime(tc.input))
		})
	}
}

func TestTimeToDailyStringFormat(t *testing.T) {
	tests := []struct {
		name   string
		input  time.Time
		output string
	}{
		{
			name:   "exist",
			input:  time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC),
			output: fmt.Sprintf("%d%02d%02d", time.Now().Year(), time.Now().Month(), time.Now().Day()),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.output, TimeToDailyStringFormat(tc.input))
		})
	}
}

func TestHourlyStringToTime(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output time.Time
	}{
		{
			name:   "empty",
			output: time.Time{},
		},
		{
			name:   "exist",
			input:  fmt.Sprintf("%d%02d%02d%02d", time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour()),
			output: time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.output, HourlyStringToTime(tc.input))
		})
	}
}

func TestTimeToHourlyStringFormat(t *testing.T) {
	tests := []struct {
		name   string
		input  time.Time
		output string
	}{
		{
			name:   "exist",
			input:  time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), 0, 0, 0, time.UTC),
			output: fmt.Sprintf("%d%02d%02d%02d", time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour()),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.output, TimeToHourlyStringFormat(tc.input))
		})
	}
}
