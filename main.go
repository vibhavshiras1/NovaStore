package main

import (
	"fmt"
	"log"

	"github.com/vibhavshiras1/novastore/p2p"
)

func OnPeer(peer p2p.Peer) error {
	peer.Close()
	return nil
}

func main() {
	fmt.Println("Distributed file systems storage")

	addr := ":3000"

	tcpOpts := p2p.TCPTransportOpts{
		ListenAddress: addr,
		HandShakeFunc: p2p.NOPHandShakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
	}

	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()

	err := tr.ListenAndAccept()

	if err != nil {
		log.Fatal(err)
	}

	select {}
}
