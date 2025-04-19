package p2p

import "net"

//remote node
type Peer interface{
	Send([]byte) error
	net.Conn
}


//anything for communicattion between nodes
type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}