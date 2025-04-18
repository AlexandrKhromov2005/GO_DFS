package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)  

func TestPathTransformFunc(t *testing.T) {
	key := "test_key"
	pathKey := CASPathTransformFunc(key)
	expectedOriginalKey := "00942f4668670f34c5943cf52c7ef3139fe2b8d6"
	expectedPathName := "00942/f4668/670f3/4c594/3cf52/c7ef3/139fe/2b8d6"

	if pathKey.PathName != expectedPathName {
		t.Errorf("Expected %s, got %s", expectedPathName, expectedPathName)
	}

	if pathKey.Original != expectedOriginalKey {
		t.Errorf("Expected %s, got %s", pathKey.Original, expectedOriginalKey)
	}
}

func TestStore(t *testing.T) {
	s := newStore()
	defer teardown(t, s)

	for i := 0; i < 50; i++ {

		key := fmt.Sprintf("foo_%d", i)
		data := []byte("some jpg bytes")

		if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); !ok {
			t.Errorf("Expected key %s to exist", key)
		}

		r, err := s.Read(key)
		if err != nil {
			t.Error(err)
		}

		b, _ := io.ReadAll(r)
		if !bytes.Equal(b, data) {
			t.Errorf("Expected %s, got %s", data, b)
		}

		fmt.Println(string(b))

		if err := s.Delete(key); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); ok {
			t.Errorf("Expected key %s to not exist", key)
		}
	}
} 

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	return NewStore(opts)
}

func teardown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}