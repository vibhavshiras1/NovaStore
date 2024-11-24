package p2p

// Peer is a remote node
type Peer interface {
}

// Transport handles communication between peers (nodes)
type Transport interface {
	ListenAndAccept() error
}
