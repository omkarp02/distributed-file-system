package main

import (
	"fmt"
	"log"

	"github.com/omkarp02/distributed-file-system/p2p"
)

func OnPerr(peer p2p.Peer) error {
	fmt.Println("doing some logic with the peer outside of tcptransport")
	peer.Close()
	return nil
}

func main() {

	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandeshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPerr,
	}

	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			log.Println(msg)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal()
	}

	select {}
}
