package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/vibhavshiras1/novastore/p2p"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BootStrapNodes    []string
}

type FileServer struct {
	FileServerOpts

	peerLock sync.Mutex
	peers    map[string]p2p.Peer
	store    *Store
	quitch   chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}
	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitch:         make(chan struct{}),
		peers:          make(map[string]p2p.Peer),
	}
}

type Message struct {
	Payload any
}

func (s *FileServer) broadcast(msg *Message) error {
	peers := []io.Writer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)

	return gob.NewEncoder(mw).Encode(msg)
}

func (s *FileServer) StoreData(key string, r io.Reader) error {

	buf := new(bytes.Buffer)
	msg := Message{
		Payload: []byte("storagekey"),
	}

	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		fmt.Print(err)
		return err
	}

	for _, peer := range s.peers {
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}

	time.Sleep(time.Second * 3)

	payload := []byte("THIS LARGE FILE")
	for _, peer := range s.peers {
		if err := peer.Send(payload); err != nil {
			return err
		}
	}

	return nil

	// buf := new(bytes.Buffer)
	// tee := io.TeeReader(r, buf)

	// if err := s.store.Write(key, tee); err != nil {
	// 	return err
	// }

	// p := &DataMessage{
	// 	Key:  key,
	// 	Data: buf.Bytes(),
	// }

	// // fmt.Println(buf.Bytes())
	// return s.broadcast(&Message{
	// 	From:    "todo",
	// 	payload: p,
	// })
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[p.RemoteAddr().String()] = p

	log.Printf("Connected with remote: %s", p.RemoteAddr().String())
	return nil
}

func (s *FileServer) loop() {

	defer func() {
		log.Println("File Server stopped due to user quit action")
		s.Transport.Close()
	}()

	for {
		select {
		case rpc := <-s.Transport.Consume():
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Recv: %s\n", string(msg.Payload.([]byte)))

			peer, ok := s.peers[rpc.From]
			if !ok {
				panic("Peer not found in peers map")
			}

			b := make([]byte, 1000)
			if _, err := peer.Read(b); err != nil {
				panic(err)
			}

			fmt.Printf("%s\n", string(b))

			peer.(*p2p.TCPPeer).Wg.Done()
			// fmt.Printf("%+v\n", string(p.Data))
			// if err := s.handleMessage(&m); err != nil {
			// 	log.Println(err)
			// }

		case <-s.quitch:
			return
		}
	}
}

// func (s *FileServer) handleMessage(msg *Message) error {
// 	switch v := msg.payload.(type) {
// 	case *DataMessage:
// 		fmt.Printf("Received Data: %+v\n", v)
// 	}

// 	return nil
// }

func (s *FileServer) bootStrapNetwork() error {
	for _, addr := range s.BootStrapNodes {
		if len(addr) == 0 {
			continue
		}
		log.Println("Attempting to connect with remote: ", addr)
		go func(addr string) {
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("Dial error ", err)
			}
		}(addr)
	}

	return nil
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	if len(s.BootStrapNodes) != 0 {
		s.bootStrapNetwork()
	}

	s.loop()

	return nil
}
