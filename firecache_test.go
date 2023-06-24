package firecache

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"reflect"
	"testing"
)

var (
	firecache *Firecache

	testCollection          string
	testDocument            string
	testDocumentNotExisting string

	testData          *Document
	testDataForUpdate U
	testDataUpdated   *Document
)

func TestMain(m *testing.M) {
	firecache = GetFirecache()

	u := make([]byte, 16)
	rand.Read(u)

	u[8] = (u[8] | 0x80) & 0xBF
	u[6] = (u[6] | 0x40) & 0x4F

	testCollection = "suites/go-firecache/test-" + hex.EncodeToString(u)
	testDocument = "test_doc"
	testDocumentNotExisting = "test_doc_not_existing"

	testData = &Document{"test_field": "test_value"}
	testDataForUpdate = U{{Path: "test_field", Value: "test_value_updated"}}
	testDataUpdated = &Document{"test_field": "test_value_updated"}

	code := m.Run()

	os.Exit(code)
}

func TestInsert(t *testing.T) {
	id, err := firecache.Insert(testCollection, testData)
	if err != nil {
		t.Error(err)
	}
	if id == "" {
		t.Errorf("not existing document could not be added to collection")
	}

	_, err = firecache.Insert(testCollection+"/"+testDocument, testData)
	if err != nil {
		t.Error(err)
	}

	_, err = firecache.Insert(testCollection+"/"+testDocument, testData)
	if err == nil {
		t.Errorf("already existing document was not rejected by the method")
	}
}

func TestRead(t *testing.T) {
	_, err := firecache.Read(testCollection+"/"+testDocumentNotExisting, nil)
	if err == nil {
		t.Errorf("not existing document was found by the method")
	}

	doc, _ := firecache.Read(testCollection+"/"+testDocument, nil)

	if !reflect.DeepEqual(doc, testData) {
		t.Errorf("expected %v, got %v", testData, doc)
	}
}

func TestUpdate(t *testing.T) {
	err := firecache.Update(testCollection+"/"+testDocumentNotExisting, testDataForUpdate)
	if err == nil {
		t.Errorf("not existing document was updated by the method")
	}

	err = firecache.Update(testCollection, testDataForUpdate)
	if err == nil {
		t.Errorf("collection path was not rejected by the method")
	}

	err = firecache.Update(testCollection+"/"+testDocument, testDataForUpdate)
	if err != nil {
		t.Error(err)
	}

	doc, _ := firecache.ReadWithoutCache(testCollection+"/"+testDocument, nil)
	if !reflect.DeepEqual(doc, testDataUpdated) {
		t.Errorf("expected %v, got %v", testDataUpdated, doc)
	}
}

func TestDelete(t *testing.T) {
	err := firecache.Delete(testCollection+"/"+testDocumentNotExisting, nil)
	if err == nil {
		t.Errorf("not existing document was deletable in the method")
	}

	err = firecache.Delete(testCollection+"/"+testDocument, nil)
	if err != nil {
		t.Error("document has not been deleted")
	}

	_, err = firecache.ReadWithoutCache(testCollection+"/"+testDocument, nil)
	if err == nil {
		t.Errorf("deleted document was found by the method")
	}

	err = firecache.Delete(testCollection, nil)
	if err != nil {
		t.Error("collection has not been deleted")
	}

	res, _ := firecache.Read(testCollection, nil)
	docs := ParseColl(res)
	if l := len(*docs); l > 0 {
		t.Errorf("expected 0 entries, got %v", l)
	}
}
