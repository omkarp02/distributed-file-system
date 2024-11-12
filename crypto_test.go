package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCopyEncrypt(t *testing.T) {
	payload := "Foo not bar"
	src := bytes.NewReader([]byte(payload))
	dst := new(bytes.Buffer)
	key := newEncryptionKey()
	_, err := copyEncrypt(key, src, dst)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(dst.String())

	out := new(bytes.Buffer)
	nw, err := copyDecrypt(key, dst, out)
	if err != nil {
		t.Error(err)
	}

	if nw != 16+len(payload) {

	}

	fmt.Println(out.String())

}

func TestNewEncryptionKey(t *testing.T) {
	key := newEncryptionKey()
	fmt.Println(key)
}
