package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listen_adress := ":4000"
	tr := NewTCPTransport(listen_adress)

	assert.Equal(t, tr.listenAddress, listen_adress)

	err := tr.ListenAndAccept()
	assert.Nil(t, err)
}
