package dal

import (
	"reflect"
	"testing"
	"time"
)

func TestGetPurchase(t *testing.T) {
	db, err := GetMockupDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	purchases, err := GetPurchases(db)
	if err != nil {
		t.Error(err)
	}

	if len(purchases) != 2 {
		t.Error("expected 2 purchases")
	}

	mockup1 := Purchase{
		uint64(1), "Test1", 20.0, time.Unix(0, 0),
	}
	if !reflect.DeepEqual(purchases[0], mockup1) {
		t.Error("first retrieved row should match expected")
	}

	mockup2 := Purchase{
		uint64(2), "Test2", 30.0, time.Unix(0, 0),
	}
	if !reflect.DeepEqual(purchases[1], mockup2) {
		t.Error("second retrieved row should match expected")
	}
}
