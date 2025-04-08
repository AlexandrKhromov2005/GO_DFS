package p2p

//any arbitrary data that can be sent over the each transport between two peeers
type Message struct {
	Payload []byte

}