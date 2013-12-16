package rockgo

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

func TestOpen(t *testing.T) {
	db := New()

	dbName := tempDBPath()

	// open db w/ createIfMissing = false, should produce error
	err := db.Open(dbName, false)
	if err == nil {
		t.Errorf("Open should return error when opening non-existing db")
	}

	// open should create db if createIfMissing = true
	err = db.Open(dbName, true)
	if err != nil {
		t.Errorf("Open failed: %v", err)
	}
	deleteDB(t, dbName)
}

func TestPutGet(t *testing.T) {
	db, dbPath := initDB(t)
	defer deleteDB(t, dbPath)

	k := []byte("foo")
	v := []byte("bar")

	err := db.Put(k, v)
	if err != nil {
		t.Errorf("Put failed: %v", err)
	}

	v2, err := db.Get(k)
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if string(v2) != string(v) {
		t.Errorf("Expected Get to return %v, got %v", v, v2)
	}
}

func TestDelete(t *testing.T) {
	db, dbPath := initDB(t)
	defer deleteDB(t, dbPath)

	k := []byte("foo")
	v := []byte("bar")

	err := db.Put(k, v)
	if err != nil {
		t.Errorf("Put failed: %v", err)
	}

	err = db.Delete(k)
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	v2, err := db.Get(v)
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if v2 != nil {
		t.Errorf("Expected Get of previously deleted key to return nil, got %v", v2)
	}
}

func tempDBPath() string {
	dbName := fmt.Sprintf("rockgo-test-%d", rand.Int())
	path := filepath.Join(os.TempDir(), dbName)
	return path
}

func initDB(t *testing.T) (db *DB, dbPath string) {
	db = New()
	dbPath = tempDBPath()
	err := db.Open(dbPath, true)

	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	return
}

func deleteDB(t *testing.T, path string) {
	err := os.RemoveAll(path)
	if err != nil {
		t.Errorf("Unable to remove database dir: %s", path)
	}
}
