package dal

import "time"

const marker = 24

func currentMarkerDay() time.Time {
	return currentMarkerFromDay(time.Now())
}

func currentMarkerFromDay(base time.Time) time.Time {
	thisMonth := base.Month()
	thisYear := base.Year()

	markerProto := time.Date(thisYear, thisMonth, marker, 0, 0, 0, 0, time.Local)
	if base.Day() <= marker {
		return markerProto.AddDate(0, -1, 0)
	}

	return markerProto
}
