package p2p

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	opts := TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: NOPHandshakeFunc,
		Decoder:       &GOBDecoder{},
	}
	tr := NewTCPTransport(opts)

	assert.Equal(t, opts.ListenAddr,":3000")

	assert.Nil(t, tr.ListenAndAccept())

}