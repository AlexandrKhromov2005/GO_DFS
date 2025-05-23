package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
)

//TCPPeer represents the remote node over a TCP established connection
type TCPPeer struct {
	//the underlying connecting of the peer, witch in this case
	//is a TCP connection
	net.Conn
	//if we dial and recieve a conn => outbound == true
	//if we accept and retrieve a conn => outbound == false
	outbound 	bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn: 		conn,
		outbound: 	outbound,
	}
}

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.Conn.Write(b)
	return err
}

//Dial implements Transport interface
func (t *TCPTransport) Dial (addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go t.handleConn(conn, true)

	return nil
}

type TCPTransportOpts struct {
	ListenAddr		string
	HandshakeFunc	HandshakeFunc
	Decoder			Decoder
	OnPeer			func(Peer) error
}
 
type TCPTransport struct {
	TCPTransportOpts
	listener		net.Listener
	rpcch 			chan RPC
}


func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: 	opts,
		rpcch: 				make(chan RPC),
	}
}

//Consume implements the transport interface,
//which will return read only channel for reading 
//the incomming messages recieved from another peer in th network
func (t *TCPTransport) Consume() <- chan RPC {
	return t.rpcch
}

//Close implements Transport interface
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener , err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	log.Printf("TCP transport listening on port : %s\n", t.ListenAddr)

	return nil
}

func (t *TCPTransport) startAcceptLoop()  {
	for {
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

type Temp struct {

}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool){
	var err error
	defer func(){
		fmt.Printf("dropping peer connection: %s\n", err)
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

	//Read loop
	rpc := RPC{}
	for {
		err = t.Decoder.Decode(conn, &rpc)
		if err != nil {
			return 
		}

		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc

	}
}