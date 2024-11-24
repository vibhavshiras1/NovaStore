package p2p

import "errors"

var ErrInvalidHandshake = errors.New("Invalid Hanshake")

type HandShakeFunc func(Peer) error

func NOPHandShakeFunc(Peer) error { return nil }
