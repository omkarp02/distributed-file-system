package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	// key := "momsbestpictrue"
	// pathname := CASPathTransformFunc(key)

	// fmt.Println(pathname)

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

func TestDelete(t *testing.T) {
	s := newStore()

	key := "fooandbar"

	data := []byte("some jpg bytes")

	if _, err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}

func TestStore(t *testing.T) {
	s := newStore()

	defer teardown(t, s)

	key := "momsspecials"

	data := []byte("some jpg bytes")

	if _, err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if ok := s.Has(key); !ok {
		t.Errorf("expected to have key %s", key)
	}

	_, r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := io.ReadAll(r)

	fmt.Println("no of bytes written to disc", b)

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}

	if ok := s.Has(key); ok {
		t.Errorf("key exist meaning delte not working")
	}

}
