package dal

import (
	"reflect"
	"testing"
	"time"
)

func TestGetPurchases(t *testing.T) {
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
		uint64(1), "Test1", 20.0, time.Unix(1444885951, 0),
	}
	if !reflect.DeepEqual(purchases[0], mockup1) {
		t.Error("first retrieved row should match expected", purchases[0], mockup1)
	}

	mockup2 := Purchase{
		uint64(2), "Test2", 30.0, time.Unix(1444885950, 0),
	}
	if !reflect.DeepEqual(purchases[1], mockup2) {
		t.Error("second retrieved row should match expected")
	}
}

func TestAddPurchase(t *testing.T) {
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

	newPurchase := Purchase{
		uint64(3), "Rawr", 332.03, time.Now(),
	}
	err = AddPurchase(db, &newPurchase)
	if err != nil {
		t.Error(err)
	}

	purchases, err = GetPurchases(db)
	if err != nil {
		t.Error(err)
	}
	if len(purchases) != 3 {
		t.Error("expected 3 purchases now, AddPurchase failed")
	}
	if newPurchase.Name != purchases[0].Name && newPurchase.Cost != purchases[0].Cost {
		t.Error("new purchase was not saved correctly", newPurchase, purchases[0])
	}
}

func TestGetPurchase(t *testing.T) {
	db, err := GetMockupDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mockup1 := Purchase{
		uint64(1), "Test1", 20.0, time.Unix(1444885951, 0),
	}
	purchase1, err := GetPurchase(db, 1)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(*purchase1, mockup1) {
		t.Error("first retrieved row should match expected")
	}

	mockup2 := Purchase{
		uint64(2), "Test2", 30.0, time.Unix(1444885950, 0),
	}
	purchase2, err := GetPurchase(db, 2)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(*purchase2, mockup2) {
		t.Error("second retrieved row should match expected")
	}
}

func TestDeletePurchase(t *testing.T) {
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
		t.Error("expected 2 purchases initially")
	}

	err = DeletePurchase(db, 1)
	if err != nil {
		t.Error(err)
	}

	purchases, err = GetPurchases(db)
	if err != nil {
		t.Error(err)
	}
	if len(purchases) != 1 {
		t.Error("expected 1 purchase after delete")
	}
	mockup2 := Purchase{
		uint64(2), "Test2", 30.0, time.Unix(1444885950, 0),
	}
	if !reflect.DeepEqual(purchases[0], mockup2) {
		t.Error("second retrieved row should match expected")
	}

}
