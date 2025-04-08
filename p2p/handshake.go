package p2p

import "errors"


//returns if handshake between remote and local node could not be established
var ErrInvalidHandshake = errors.New("invalid handshake")

type HandshakeFunc func(Peer) error

func NOPHandshakeFunc(Peer) error {
	return nil
}