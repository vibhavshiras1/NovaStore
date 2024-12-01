package p2p

import (
	"fmt"
	"net"
	"sync"
)

type TCPPeer struct {
	conn net.Conn

	// dial and retrieve -> outbound = True
	// accept and retrieve -> outbound = False
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{conn: conn, outbound: outbound}
}

func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOpts struct {
	ListenAddress string
	HandShakeFunc HandShakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan RPC

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC),
	}
}

func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddress)

	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil

}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP Accept error: %s\n", err)
		}

		fmt.Printf("New incoming connection: %+v\n", conn)
		go t.handleConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {
	var err error
	defer func() {
		fmt.Printf("Dropping peer connection: %s", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)

	err = t.HandShakeFunc(peer)
	if err != nil {
		return
	}

	if t.OnPeer != nil {
		err = t.OnPeer(peer)
		if err != nil {
			return
		}
	}

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
