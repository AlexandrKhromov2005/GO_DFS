package main

import (
	"bytes"
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
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	s := NewStore(opts)
	key := "test_key"
	data := []byte("some jpg bytes")

	if err := s.writeStream("test_key", bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := io.ReadAll(r)

	if !bytes.Equal(b, data) {
		t.Errorf("Expected %s, got %s", data, b)
	}	
}