package goa

import (
	"net"
	"sync"
)

// AudioServer contain output from audio server
type AudioServer struct {
	pc     net.PacketConn
	buffer chan []byte
	wg     sync.WaitGroup
}

// AudioConn hold udp connection to the server
type AudioConn struct {
	conn     *net.UDPConn
	bufferIn chan []byte
	header   []byte
}
