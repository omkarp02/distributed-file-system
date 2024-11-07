package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/omkarp02/distributed-file-system/p2p"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BoostrapNodes     []string
}

type FileServer struct {
	FileServerOpts

	peerLock sync.Mutex
	peers    map[string]p2p.Peer

	store  *Store
	quitch chan struct{}
}

type Message struct {
	From    string
	Payload any
}

type DataMessage struct {
	Key  string
	Data []byte
}

func init() {
	gob.Register(DataMessage{})
}

func NewFileServer(opts FileServerOpts) *FileServer {

	storeOpts := StoreOpts{
		PathTransformFunc: opts.PathTransformFunc,
		Root:              opts.StorageRoot,
	}

	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitch:         make(chan struct{}),
		peers:          make(map[string]p2p.Peer),
	}
}

func (s *FileServer) StoreData(key string, r io.Reader) error {
	//1. Store this file to disk
	//2. broadcast this file to all known peers in the network

	buf := new(bytes.Buffer)

	tee := io.TeeReader(r, buf)

	if err := s.store.Write(key, tee); err != nil {
		fmt.Println(err)
		return err
	}

	dataMsg := &DataMessage{
		Key:  key,
		Data: buf.Bytes(),
	}

	// the reader is empty

	return s.broadcast(&Message{
		From:    "todo",
		Payload: dataMsg,
	})
}

func (s *FileServer) broadcast(msg *Message) error {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(">>>>>>")

	for _, peer := range s.peers {
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}

	return nil

}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[p.RemoteAddr().String()] = p

	log.Printf("connected with remote %s", p.RemoteAddr())

	return nil

}

func (s *FileServer) loop() {

	defer func() {
		log.Println("file server stopped due touser quit action")
		s.Transport.Close()
	}()

	for {
		select {
		case rpc := <-s.Transport.Consume():
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Fatal(">>>>>>> err", err)
			}

			log.Printf("recv: %s", msg)

			// if err := s.handleMessage(&msg); err != nil {
			// 	log.Println(err)
			// }

		case <-s.quitch:
			return
		}
	}
}

func (s *FileServer) handleMessage(msg *Message) error {
	switch v := msg.Payload.(type) {
	case *DataMessage:
		log.Printf("recieved data %v+\n", v)
	}

	return nil
}

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BoostrapNodes {

		if len(addr) == 0 {
			continue
		}

		go func(addr string) {
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("dial error:", err)
			}
		}(addr)
	}

	return nil
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	s.bootstrapNetwork()

	s.loop()

	return nil
}
