package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCopyEncrypt(t *testing.T) {
	src := bytes.NewReader([]byte("foo not bar"))
	dst := new(bytes.Buffer)
	key := newEncryptionKey()
	_, err := copyEncrypt(key, src, dst)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(dst.String())

	out := new(bytes.Buffer)
	if _, err := copyDecrypt(key, dst, out); err != nil {
		t.Error(err)
	}

	fmt.Println(out.String())

}

func TestNewEncryptionKey(t *testing.T) {
	key := newEncryptionKey()
	fmt.Println(key)
}
