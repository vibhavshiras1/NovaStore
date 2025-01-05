package p2p

import "net"

// Peer is a remote node
type Peer interface {
	net.Conn
	Send([]byte) error
	// RemoteAddr() net.Addr
	// Close() error
	CloseStream()
}

// Transport handles communication between peers (nodes)
type Transport interface {
	Addr() string
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
