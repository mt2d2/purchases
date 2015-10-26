package dal

import (
	"reflect"
	"testing"
	"time"
)

func TestGetPurchases(t *testing.T) {
	db, err := mockupDB()
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
	if !reflect.DeepEqual(purchases[0], *mockup1()) {
		t.Error("first retrieved row should match expected", purchases[0], *mockup1())
	}
	if !reflect.DeepEqual(purchases[1], *mockup2()) {
		t.Error("second retrieved row should match expected")
	}
}

func TestGetPurchasesAfterDate(t *testing.T) {
	db, err := mockupDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	beginningOfTime := time.Unix(mockup1().TimeBought.Unix()-(60*60*24*1), 0)
	purchases, err := getPurchasesAfterDate(db, beginningOfTime)
	if err != nil {
		t.Error(err)
	}
	if len(purchases) != 2 {
		t.Error("expected 2 purchases ")
	}

	justBefore := mockup1().TimeBought
	purchases, err = getPurchasesAfterDate(db, justBefore)
	if err != nil {
		t.Error(err)
	}
	if len(purchases) != 1 {
		t.Error("expected 1 purchases", purchases)
	}

	now := time.Now()
	purchases, err = getPurchasesAfterDate(db, now)
	if err != nil {
		t.Error(err)
	}
	if len(purchases) != 0 {
		t.Error("expected 0 purchases")
	}
}

func TestValidatePurchase(t *testing.T) {
	test := mockup1()

	ok, errs := ValidatePurchase(test)
	if !ok || len(errs) != 0 {
		t.Error("expected a valid test")
	}

	test.Cost = 0
	ok, errs = ValidatePurchase(test)
	if ok || len(errs) != 1 {
		t.Error("expected a 1 error")
	}

	test.Name = ""
	ok, errs = ValidatePurchase(test)
	if ok || len(errs) != 2 {
		t.Error("expected a 2 errors")
	}

	oneHour, _ := time.ParseDuration("1h")
	test.TimeBought = time.Now().Add(oneHour)
	ok, errs = ValidatePurchase(test)
	if ok || len(errs) != 3 {
		t.Error("expected a 3 errors")
	}
}

func TestAddPurchase(t *testing.T) {
	db, err := mockupDB()
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
	if newPurchase.Name != purchases[0].Name &&
		newPurchase.Cost != purchases[0].Cost &&
		!newPurchase.TimeBought.Before(time.Now()) {
		t.Error("new purchase was not saved correctly", newPurchase, purchases[0])
	}
}

func TestGetPurchase(t *testing.T) {
	db, err := mockupDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	purchase1, err := GetPurchase(db, 1)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(*purchase1, *mockup1()) {
		t.Error("first retrieved row should match expected")
	}

	purchase2, err := GetPurchase(db, 2)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(*purchase2, *mockup2()) {
		t.Error("second retrieved row should match expected")
	}
}

func TestDeletePurchase(t *testing.T) {
	db, err := mockupDB()
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
	if !reflect.DeepEqual(purchases[0], *mockup2()) {
		t.Error("second retrieved row should match expected")
	}

}
