package p2p

import "net"

//any arbitrary data that can be sent over the each transport between two peeers
type RPC struct {
	From net.Addr
	Payload []byte

}