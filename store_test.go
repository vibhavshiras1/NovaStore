package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	return NewStore(opts)
}

func teardown(t *testing.T, store *Store) {
	if err := store.Clear(); err != nil {
		t.Error(err)
	}
}

func TestTransformPathFunc(t *testing.T) {
	key := "greatestintheworld"
	pathKey := CASPathTransformFunc(key)

	expectedFileName := "abb9f44788eef303f9290238437fe418cd166c8e"
	expectedPathName := "abb9f/44788/eef30/3f929/02384/37fe4/18cd1/66c8e"
	if pathKey.pathName != expectedPathName {
		t.Errorf("Actual: %s, Expected: %s", pathKey.pathName, expectedPathName)
	}
	if pathKey.FileName != expectedFileName {
		t.Errorf("Actual: %s, Expected: %s", pathKey.pathName, expectedFileName)
	}
}

func TestStore(t *testing.T) {
	store := newStore()
	defer teardown(t, store)

	key := "lowlevellearning"
	data := []byte("Some jpg file")

	// Writing Data
	err := store.Write(key, bytes.NewReader(data))
	if err != nil {
		t.Error(err)
	}

	if exists := store.Has(key); !exists {
		t.Errorf("Key: %s not found", key)
	}
	fmt.Printf("Key: %s exists\n", key)

	// Reading Data
	r, err := store.Read(key)
	b, _ := io.ReadAll(r)

	if string(b) != string(data) {
		t.Errorf("Actual: %s, Expected: %s", string(b), string(data))
	}
	fmt.Printf("Key: %s, Data: %s\n", key, string(b))

	// Deleting the key and its contents
	del_err := store.Delete(key)
	if del_err != nil {
		t.Error(del_err)
	}
	fmt.Printf("Successfully deleted the key: %s", key)
}

func TestStoreMultipleKeys(t *testing.T) {
	store := newStore()
	defer teardown(t, store)

	keyArr := make([]string, 0, 10)
	for i := 0; i < 10; i++ {
		keyArr = append(keyArr, fmt.Sprintf("newkey_%d", i))
	}
	data := []byte("Distributed file systems is amazing")

	// Writing Multiple Keys
	for i := 0; i < 10; i++ {
		key := keyArr[i]

		err := store.Write(key, bytes.NewReader(data))
		if err != nil {
			t.Error(err)
		}
	}

	// Checking if all the keys exist
	for i := 0; i < 10; i++ {
		key := keyArr[i]
		if exists := store.Has(key); !exists {
			t.Errorf("Key: %s not found", key)
		}
		fmt.Printf("Key: %s exists\n", key)
	}

	// Reading all the keys
	for i := 0; i < 10; i++ {
		key := keyArr[i]
		r, _ := store.Read(key)
		b, _ := io.ReadAll(r)

		if string(b) != string(data) {
			t.Errorf("Actual: %s, Expected: %s", string(b), string(data))
		}
		fmt.Printf("Key: %s, Data: %s\n", key, string(b))
	}

	// Deleting all the keys
	for i := 0; i < 10; i++ {
		key := keyArr[i]
		del_err := store.Delete(key)
		if del_err != nil {
			t.Error(del_err)
		}
		fmt.Printf("Successfully deleted the key: %s", key)
	}

	// Verifying that the keys don't exist anymore
	for i := 0; i < 10; i++ {
		key := keyArr[i]
		if exists := store.Has(key); exists {
			t.Errorf("Key: %s still exists", key)
		}
	}
}
