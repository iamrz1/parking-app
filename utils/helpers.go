package utils

import "time"

const (
	ChargePerHour = 10
)

func BoolP(in bool) *bool {
	return &in
}

func DurationToHours(d time.Duration) int {
	const hour = time.Hour

	hours := int(d / hour)

	if d%hour > 0 {
		hours++
	}

	return hours
}
