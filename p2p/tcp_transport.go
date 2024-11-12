package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

// TCPpeer represents the remote node over a tcp established connection.
type TCPPeer struct {
	// The underlying connection of the peer which in this
	// case is tcp connection.
	//so tcppeer can directly use all the method of net.Conn
	//like tcppeer.read()
	net.Conn
	// if we dial and retrieve a conn => outbound == true
	// if we accept and retrieve a conn => outbound == false
	outbound bool

	wg *sync.WaitGroup
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
		wg:       &sync.WaitGroup{},
	}
}

func (p *TCPPeer) CloseStream() {
	p.wg.Done()
}

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.Conn.Write(b)
	return err
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan RPC
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC, 1024),
	}
}

// Consume implements the transport interface, which will return read only channel
// for reading the incoming msg received from another peer
func (t *TCPTransport) Consume() <-chan RPC {

	return t.rpcch
}

// addr implement Transport interface return the address
// the transport is accepting
func (t *TCPTransport) Addr() string {
	return t.ListenAddr
}

// close implement Transport interface
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

// Dial implement Transport interface
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go t.handleConn(conn, true)

	return nil
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	log.Println("tcp transport listening on port: ", t.ListenAddr)

	return nil

}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()

		if errors.Is(err, net.ErrClosed) {
			return
		}

		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}

		fmt.Printf("new incoming connection %+v\n", conn)
		go t.handleConn(conn, false)

	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {

	var err error

	defer func() {
		fmt.Println("dropping peer connectin", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, outbound)

	if err := t.HandshakeFunc(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			return
		}
	}

	for {
		//Read Loop
		rpc := RPC{}
		err = t.Decoder.Decode(conn, &rpc)
		fmt.Printf("someone send message to this tcp connection\ns")

		//TODO: here we need to figure out if there is a conn close error then only return it, as now we are returning for every error
		if err != nil {
			fmt.Printf("TCP read error: %s\n", err)
			return
		}

		rpc.From = conn.RemoteAddr().String()

		if rpc.Stream {
			peer.wg.Add(1)
			fmt.Printf("[%s] incoming stream, waiting...\n", conn.RemoteAddr())
			peer.wg.Wait()
			fmt.Printf("[%s] stream closed, resuming read loop\n", conn.RemoteAddr())
			continue
		}

		t.rpcch <- rpc
	}

}
