package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/omkarp02/distributed-file-system/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandeshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}

	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := FileServerOpts{
		EncKey:            newEncryptionKey(),
		StorageRoot:       listenAddr[1:] + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BoostrapNodes:     nodes,
	}

	s := NewFileServer(fileServerOpts)

	tcpTransport.OnPeer = s.OnPeer

	return s

}

func main() {

	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", "")
	s3 := makeServer(":5000", ":3000", ":4000")

	go func() {
		log.Fatal(s1.Start())
	}()
	time.Sleep(time.Millisecond * 500)
	go func() {
		log.Fatal(s2.Start())
	}()

	time.Sleep(2 * time.Second)

	go s3.Start()
	time.Sleep(2 * time.Second)

	for i := 0; i < 3; i++ {
		// key := "coolPicture"
		key := fmt.Sprintf("picture_%d.png", i)

		data := bytes.NewReader([]byte("my big data file here"))
		s3.Store(key, data)

		if err := s3.store.Delete(key); err != nil {
			log.Fatal(err)
		}

		r, err := s3.Get(key)
		if err != nil {
			log.Fatal("err")
		}

		defer r.Close()

		b, err := io.ReadAll(r)
		if err != nil {
			log.Fatal("err")
		}
		fmt.Println(string(b))
	}

}

//timeline 8:27:02
