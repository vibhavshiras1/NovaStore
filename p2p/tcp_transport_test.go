package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listen_adress := ":3000"
	tcpOpts := TCPTransportOpts{
		ListenAddress: listen_adress,
		HandShakeFunc: NOPHandShakeFunc,
		Decoder:       DefaultDecoder{},
	}
	tr := NewTCPTransport(tcpOpts)

	assert.Equal(t, tr.ListenAddress, listen_adress)

	err := tr.ListenAndAccept()
	assert.Nil(t, err)
}
