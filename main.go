package main

import (
	"fmt"
	"log"

	"github.com/vibhavshiras1/novastore/p2p"
)

func main() {
	fmt.Println("Distributed file systems storage")

	addr := ":3000"

	tcpOpts := p2p.TCPTransportOpts{
		ListenAddress: addr,
		HandShakeFunc: p2p.NOPHandShakeFunc,
		Decoder:       p2p.GOBDecoder{},
	}

	tr := p2p.NewTCPTransport(tcpOpts)

	err := tr.ListenAndAccept()

	if err != nil {
		log.Fatal(err)
	}

	select {}
}
