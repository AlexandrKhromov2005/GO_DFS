package p2p

//remote node
type Peer interface{
	Close() error
}


//anything for communicattion between nodes
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}