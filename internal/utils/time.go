package utils

import (
	"fmt"
	"time"
)

func FormatDate(ts time.Time) string {
	return ts.Format("01-02-06")
}

func FormatTimeSince(ts *time.Time) string {
	if ts == nil {
		return "never"
	}
	now := time.Now()
	elapsed := now.Sub(*ts).Round(time.Second)
	return fmt.Sprintf("%v ago", elapsed)
}
