package utils

import (
	"testing"
	"time"
)

func TestDurationToHours(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected int
	}{
		{
			name:     "Exactly one hour",
			duration: time.Hour,
			expected: 1,
		},
		{
			name:     "One hour and five minutes",
			duration: time.Hour*3 + 5*time.Minute,
			expected: 4,
		},
		{
			name:     "One hour and fifty-nine minutes",
			duration: time.Hour + 59*time.Minute,
			expected: 2,
		},
		{
			name:     "Zero hours",
			duration: 0,
			expected: 0,
		},
		{
			name:     "Less than one hour",
			duration: 45 * time.Minute,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DurationToHours(tt.duration); got != tt.expected {
				t.Errorf("durationToHours(%v) = %v, want %v", tt.duration, got, tt.expected)
			}
		})
	}
}
