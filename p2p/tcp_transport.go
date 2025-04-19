package p2p

import (
	"fmt"
	"log"
	"net"
	"errors"
)

//remote node a TCP connection
type TCPPeer struct {
	//the underlying connection of this peer
	net.Conn
	//if we dial and retrieve conn => outbound == true
	//if we listen and accept conn => outbound == false
	outbound bool
}


func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn : conn,
		outbound: outbound,
	}
}

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.Conn.Write(b)
	return err
}

type TCPTransportOpts struct {
	ListenAddr string
	HandshakeFunc HandshakeFunc
	Decoder	Decoder
	OnPeer func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener      net.Listener
	shakeHands HandshakeFunc
	rpcch chan RPC       
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch: make(chan RPC, 100),
	}
}


//Consume implents the Transport interface
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

//Close implements the Transport interface
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

//Dial implements the Transport interface
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go t.handleConn(conn, true)
	return nil
}

func (t *TCPTransport) ListenAndAccept() error{
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)

	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	log.Printf("TCP transport listening on %s\n", t.ListenAddr)
	
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for{
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}
		go t.handleConn(conn, false)		
	}
}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool)  {
	var err error
	defer func(){
		fmt.Printf("dropping peer connection %+v\n", err)
		conn.Close()
		}()

	peer := NewTCPPeer(conn, outbound)

	if err = t.HandshakeFunc(peer); err != nil {
		return 
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	//read loop
	rpc := RPC{}
	for {
		err := t.Decoder.Decode(conn, &rpc)

		if err != nil {
			return
		}

		rpc.From = conn.RemoteAddr()
		log.Printf("Forwarding RPC from %s to channel\n", rpc.From)
		t.rpcch <- rpc
	}
}