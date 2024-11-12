package main

import (
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
	s2 := makeServer(":4000", ":3000")

	go func() {
		log.Fatal(s1.Start())
	}()

	time.Sleep(1 * time.Second)

	go s2.Start()
	time.Sleep(1 * time.Second)

	// data := bytes.NewReader([]byte("my big data file here"))
	// s2.Store(fmt.Sprintf("coolpicture.jpg"), data)
	// time.Sleep(5 * time.Millisecond)

	r, err := s2.Get("coolpicture.jpg")
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

//timeline 8:27:02
