package p2p

import (
	"fmt"
	"net"
	"sync"
)

//remote node a TCP connection
type TCPPeer struct {
	conn net.Conn

	//if we dial and retrieve conn => outbound == true
	//if we listen and accept conn => outbound == false
	outbound bool
}


func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn : conn,
		outbound: outbound,
	}
}

type TCPTransportOpts struct {
	ListenAddr string
	HandshakeFunc HandshakeFunc
	Decoder	Decoder
}

type TCPTransport struct {
	TCPTransportOpts
	listener      net.Listener
	shakeHands HandshakeFunc
	decoder Decoder

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
	}

}

func (t *TCPTransport) ListenAndAccept() error{
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)

	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for{
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}
		fmt.Printf("new incoming connection %+v\n", conn)
		go t.handleConn(conn)		
	}


}

func (t *TCPTransport) handleConn(conn net.Conn)  {
	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error: %s \n", err)

		return 
	}

	//read loop
	msg := &Message{}
	for {
		if err := t.Decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP error: %s \n", err)
			continue
		}

		fmt.Printf("TCP message: %+v\n", msg)
	}
}