package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {

	listenAddr := ":4000"

	tcpOpts := TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: NOPHandeshakeFunc,
		Decoder:       DefaultDecoder{},
	}

	tr := NewTCPTransport(tcpOpts)

	assert.Equal(t, tr.ListenAddr, listenAddr)

	assert.Nil(t, tr.ListenAndAccept())
}
