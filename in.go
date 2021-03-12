package goa

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/gordonklaus/portaudio"
)

// Audio struct
type audioIn struct {
	*portaudio.Stream
	buffer []float32
	rate   float64
	i      int
}

// Dial function is use to dial in audio server
// this return an audio connection
func Dial(address string) (*AudioConn, error) {
	raddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return nil, err
	}

	return &AudioConn{
		conn:     conn,
		bufferIn: make(chan []byte),
	}, nil

}

// SendAudio is a function for sending or write audio data to the server
func (aux *AudioConn) SendAudio() {
	go capture(aux.bufferIn)
	for v := range aux.bufferIn {
		_, err := aux.sendAudio(v)
		if err != nil {
			fmt.Println("error header")
			return
		}
	}
}

// Close is for close audio client connection
func (aux *AudioConn) Close() error {
	if err := aux.conn.Close(); err != nil {
		return err
	}
	return nil
}

func (aux *AudioConn) sendAudio(v []byte) (int, error) {
	mode := byte(0x81)
	mask := byte(0x80)
	frame := []byte{mode, mask | byte(len(v))}
	sd := append(frame, v...)
	n, err := aux.conn.Write(sd)
	if err != nil {
		return 0, err
	}

	err = aux.readStatus()
	if err != nil {
		return 0, err
	}
	// fmt.Println(n)
	return n, nil

}

func (aux *AudioConn) readStatus() error {
	sts := make([]byte, 1)
	_, err := aux.conn.Read(sts)
	if err != nil {
		return err
	}

	if sts[0] == 0x81 {
		return errHeaderInvalid
	}
	return nil
}

func capture(buf chan<- []byte) {
	portaudio.Initialize()
	defer portaudio.Terminate()
	aux, err := newAudio(time.Second * 1)
	if err != nil {
		return
	}
	defer aux.Close()
	err = aux.Start()
	if err != nil {
		return
	}
	for {
		aux.Read()
		b := &bytes.Buffer{}
		err := binary.Write(b, binary.BigEndian, aux.buffer)
		if err != nil {
			return
		}
		buf <- b.Bytes()
	}

	err = aux.Stop()
	if err != nil {
		return
	}

}

func newAudio(delay time.Duration) (*audioIn, error) {
	h, err := portaudio.DefaultHostApi()
	if err != nil {
		return nil, err
	}
	p := portaudio.LowLatencyParameters(h.DefaultInputDevice, h.DefaultOutputDevice)
	p.Input.Channels = 1
	p.Output.Channels = 1
	p.SampleRate = 16000
	// e := &audioIn{buffer: make([]float32, 576)}
	// fmt.Println(int(p.SampleRate * delay.Seconds()))

	e := &audioIn{buffer: make([]float32, int(p.SampleRate*delay.Seconds()))}
	e.Stream, err = portaudio.OpenStream(p, e.processAudio)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (e *audioIn) processAudio(in, out []float32) {
	for i := range out {
		out[i] = .7 * e.buffer[e.i]
		e.buffer[e.i] = in[i]
		e.i = (e.i + 1) % len(e.buffer)
	}
}
