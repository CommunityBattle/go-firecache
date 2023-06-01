package test

import (
	"reflect"
	"testing"

	firecache "github.com/CommunityBattle/go-firecache"
)

func TestListener(t *testing.T) {
	firecache := firecache.GetFirecache()

	updateSuccessfullyCaught := false

	exp := map[string]interface{}{"test_field": "test_value", "added_field": "added_value"}

	callback := func(act any) {
		if reflect.DeepEqual(exp, act) {
			updateSuccessfullyCaught = true
		}
	}

	err := firecache.AddListener("test/test_doc", nil, &callback)

	if err != nil {
		t.Errorf("could not add listener: %v", err)
	}

	err = firecache.UpdateWithoutCache("test/test_doc", exp)

	if err != nil {
		t.Errorf("could not update: %v", err)
	}

	if updateSuccessfullyCaught != true {
		t.Errorf("listener got no update")
	}
}

func TestReadWithoutCache(t *testing.T) {
	firecache := firecache.GetFirecache()

	exp := map[string]interface{}{"test_field": "test_value"}
	act, err := firecache.ReadWithoutCache("test/test_doc", nil)

	if err != nil {
		t.Errorf("could not read: %v", err)
	}

	if !reflect.DeepEqual(act, exp) {
		t.Errorf("expected %v, got %v", exp, act)
	}
}

func TestInsertWithoutCache(t *testing.T) {
	firecache := firecache.GetFirecache()

	val := map[string]interface{}{"test_field": "test_value"}
	_, err := firecache.InsertWithoutCache("test/test_inserted_doc", val)

	if err != nil {
		t.Errorf("could not insert: %v", err)
	}

	got, err := firecache.ReadWithoutCache("test/test_inserted_doc", nil)

	if err != nil {
		t.Errorf("could not read: %v", err)
	}

	if !reflect.DeepEqual(val, got) {
		t.Errorf("expected %v, got %v", val, got)
	}
}

func TestDelete(t *testing.T) {
	firecache := firecache.GetFirecache()

	err := firecache.Delete("test/test_inserted_doc", nil)

	if err != nil {
		t.Errorf("could not delete: %v", err)
	}

	_, err = firecache.ReadWithoutCache("test/test_inserted_doc", nil)

	if err == nil {
		t.Error("document has not been deleted")
	}
}
