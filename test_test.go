package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestThisShit(t *testing.T) {
	buf := make([]byte, 10)

	r := bytes.NewReader([]byte("Omkarpawar"))

	r.Read(buf)

	fmt.Println(string(buf))
}
