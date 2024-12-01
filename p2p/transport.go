package p2p

// Peer is a remote node
type Peer interface {
	Close() error
}

// Transport handles communication between peers (nodes)
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}
