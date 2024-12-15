package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/vibhavshiras1/novastore/p2p"
)

func OnPeer(peer p2p.Peer) error {
	peer.Close()
	return nil
}

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddress: listenAddr,
		HandShakeFunc: p2p.NOPHandShakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		// TODO - OnPeer func
	}

	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootStrapNodes:    nodes,
	}

	s := NewFileServer(fileServerOpts)
	tcpTransport.OnPeer = s.OnPeer

	return s
}

func main() {
	fmt.Println("Distributed File Systems Storage")
	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")

	go func() {
		log.Fatal(s1.Start())
	}()

	time.Sleep(2 * time.Second)

	go s2.Start()
	time.Sleep(2 * time.Second)

	data := bytes.NewReader([]byte("my big data file here"))
	s2.StoreData("mynewkey", data)

	select {}

	// go func() {
	// 	time.Sleep(time.Second * 3)
	// 	fileServer.Stop()
	// }()

	// if err := fileServer.Start(); err != nil {
	// 	log.Fatal(err)
	// }

	// select {}
}

// func main() {
// 	fmt.Println("Distributed file systems storage")

// 	addr := ":3000"

// 	tcpOpts := p2p.TCPTransportOpts{
// 		ListenAddress: addr,
// 		HandShakeFunc: p2p.NOPHandShakeFunc,
// 		Decoder:       p2p.DefaultDecoder{},
// 		OnPeer:        OnPeer,
// 	}

// 	tr := p2p.NewTCPTransport(tcpOpts)

// 	go func() {
// 		for {
// 			msg := <-tr.Consume()
// 			fmt.Printf("%+v\n", msg)
// 		}
// 	}()

// 	err := tr.ListenAndAccept()

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	select {}
// }
