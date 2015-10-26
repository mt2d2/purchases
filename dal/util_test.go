package dal

import (
	"testing"
	"time"
)

func TestCurrentMarkerDay(t *testing.T) {
	rightNow := time.Now()
	markerDay := currentMarkerDay()

	if markerDay.After(rightNow) {
		t.Error("marker day should be before current time")
	}
}

func TestMarkerDayFromDay(t *testing.T) {
	for i := 1; i <= 24; i++ {
		base := time.Date(2015, 8, i, 0, 0, 0, 0, time.Local)
		expected := time.Date(2015, 7, 24, 0, 0, 0, 0, time.Local)
		if !currentMarkerFromDay(base).Equal(expected) {
			t.Error("same day as marker should yield previous month marker: ")
		}
	}

	for i := 25; i <= 31; i++ {
		base := time.Date(2015, 8, i, 0, 0, 0, 0, time.Local)
		expected := time.Date(2015, 8, 24, 0, 0, 0, 0, time.Local)
		if !currentMarkerFromDay(base).Equal(expected) {
			t.Error("marker day for day after marker should still be this month")
		}
	}

	base := time.Date(2015, 9, 1, 0, 0, 0, 0, time.Local)
	expected := time.Date(2015, 8, 24, 0, 0, 0, 0, time.Local)
	if !currentMarkerFromDay(base).Equal(expected) {
		t.Error("marker day in new month should reflect last marker day")
	}
}
