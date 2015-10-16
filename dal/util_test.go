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
