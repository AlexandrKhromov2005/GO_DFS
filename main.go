package main

import (
	"bytes"
	"log"
	"strings"
	"time"

	"github.com/AlexandrKhromov2005/GO_DFS/p2p"
)

func makeServer (listenAddr string, nodes ...string) *FileServer {
	tcptransportOpts := p2p.TCPTransportOpts{
		ListenAddr: listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder: &p2p.GOBDecoder{},
		//OnPeer: nil,
 
	}

	tcptransport := p2p.NewTCPTransport(tcptransportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:      strings.ReplaceAll(listenAddr, ":", "_") + "_network",		PathTransformFunc: CASPathTransformFunc,
		Transport:        tcptransport,
		BootstrapNodes:   nodes,
	}

	s := NewFileServer(fileServerOpts)

	tcptransport.OnPeer = s.OnPeer

	return s
}	

func main() {
    s1 := makeServer(":3000", "")
    s2 := makeServer(":4000", ":3000")

    go func() {
        log.Fatal(s1.Start())
    }()


    go s2.Start()
    time.Sleep(1 * time.Second)

    data := bytes.NewReader([]byte("my big data file here!"))

    if err := s2.StoreData("myprivatedata", data); err != nil {
        log.Fatal(err)
    }

	select {}
}