package shadow

import (
	"fmt"
	"math"
	"time"
)

const (
	secondsInMinute = 60
	secondsInHour   = secondsInMinute * 60
	secondsInDay    = secondsInHour * 24
	secondsInWeek   = secondsInDay * 7
	secondsInMonth  = secondsInDay * 30
	secondsInYear   = secondsInMonth * 12
)

func DateSinceAsMessage(t time.Time) string {
	diff := time.Since(t)

	if diff.Seconds() < 1 {
		return "just now"
	}

	if diff.Seconds() < secondsInMinute {
		return fmt.Sprintf("%.f seconds ago", math.Floor(diff.Seconds()))
	}

	// minutes
	if diff.Seconds() < secondsInHour {
		if diff.Seconds() <= secondsInMinute*2 {
			return "one minutes ago"
		}

		return fmt.Sprintf("%.f minutes ago", math.Floor(diff.Minutes()))
	}

	// hours
	if diff.Seconds() < secondsInDay {
		if diff.Seconds() <= secondsInHour*2 {
			return "an hour ago"
		}

		return fmt.Sprintf("%.f hours ago", math.Floor(diff.Hours()))
	}

	// days
	if diff.Seconds() <= secondsInWeek {
		if diff.Seconds() <= secondsInDay*2 {
			return "yesterday"
		}

		return fmt.Sprintf("%.f days ago", math.Floor(diff.Seconds()/secondsInDay))
	}

	// weeks
	if diff.Seconds() <= secondsInMonth {
		if diff.Seconds() <= secondsInWeek*2 {
			return "a week ago"
		}

		return fmt.Sprintf("%.f weeks ago", math.Floor(diff.Seconds()/secondsInWeek))
	}

	// months
	if diff.Seconds() <= secondsInYear {
		if diff.Seconds() <= secondsInMonth*2 {
			return "a month ago"
		}

		return fmt.Sprintf("%.f months ago", math.Floor(diff.Seconds()/secondsInMonth))
	}

	// years
	if diff.Seconds() <= secondsInYear*2 {
		return "a year ago"
	}

	return fmt.Sprintf("%.f years ago", math.Floor(diff.Seconds()/secondsInYear))
}
