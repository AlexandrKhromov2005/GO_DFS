package p2p

import "net"

//Peer is an interface that represents the remote node
type Peer interface {
	net.Conn
	Send([]byte) error
}


//Transport is anything that handles communication between peers
type Transport interface {
	Dial(string) error
	ListenAndAccept()	error
	Consume() <- chan RPC
	Close() error
}