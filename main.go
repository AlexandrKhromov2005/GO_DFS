package main

import (
	"bytes"
	"log"
	"strings"
	"time"

	"github.com/AlexandrKhromov2005/GO_DFS/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr: listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder: p2p.DefaultDecoder{},
	}

	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       strings.ReplaceAll(listenAddr, ":", "") + "_network" ,
		PathTransformFunc: CASPathTransformFunc,
		Transport:       tcpTransport,
		BootstrapNodes: nodes,
	}
	s := NewFileServer(fileServerOpts)
	tcpTransport.OnPeer = s.OnPeer

	return s
}

func main() {
	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")


	go func() {
		log.Fatal(s1.Start())
	}()
	time.Sleep(1 * time.Second)


	go s2.Start()
	time.Sleep(1 * time.Second)

	data := bytes.NewReader([]byte("my big data file here!"))
	s2.StoreData("myprivatedata", data)

	select{}
}