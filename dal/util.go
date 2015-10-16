package dal

import "time"

const marker = 24

func currentMarkerDay() time.Time {
	rightNow := time.Now()
	thisMonth := rightNow.Month()
	thisYear := rightNow.Year()

	marker := time.Date(thisYear, thisMonth, marker, 0, 0, 0, 0, time.Local)
	return marker.AddDate(0, -1, 0)
}
