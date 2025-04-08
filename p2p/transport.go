package p2p

//remote node
type Peer interface{}


//anything for communicattion between nodes
type Transport interface {
	ListenAndAccept() error
}