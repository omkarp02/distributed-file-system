package p2p

// //errinvalidhandshake is return if the handshake betwn the local and remote node c ould not be established
// var ErrInvalidHandshake = errors.New("invalid handshake")

// HandshakeFunc... ?
type HandshakeFunc func(Peer) error

func NOPHandeshakeFunc(Peer) error { return nil }
