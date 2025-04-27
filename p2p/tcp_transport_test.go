package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	opts := TCPTransportOpts{
		ListenAddr: "localhost:3000",
		HandshakeFunc: NOPHandshakeFunc,
		Decoder: DefaultDecoder{},
	}
	tr := NewTCPTransport(opts)
	assert.Equal(t, tr.ListenAddr, "localhost:3000")

	assert.Nil(t, tr.ListenAndAccept())
}

