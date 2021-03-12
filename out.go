package goa

import (
	"bytes"
	"encoding/binary"

	"log"
	"net"
	"runtime"

	"github.com/google/logger"
	"github.com/gordonklaus/portaudio"
)

// Listen running audio server
func Listen(addr string) (*AudioServer, error) {
	lis, err := net.ListenPacket("udp", addr)
	if err != nil {
		logger.Fatalf("error resolve addr : %v\n", err)
		return nil, err
	}
	return &AudioServer{
		buffer: make(chan []byte),
		pc:     lis,
	}, nil
}

// ReadAudio is running audio server
func (pc *AudioServer) ReadAudio() {
	for i := 0; i < runtime.NumCPU(); i++ {
		go pc.readAudio()
	}
	pc.playAudio()

}

func (pc *AudioServer) readAudio() {
	for {
		// buffer := make([]byte, 2306)
		// buffer := make([]byte, 64002)
		buffer := make([]byte, 65000)
		_, addr, err := pc.pc.ReadFrom(buffer)
		if err != nil {
			logger.Infof("error read data: %v", err)
			return
		}
		err = pc.readHeader(buffer, addr)
		if err != nil {
			log.Println(err)
			return
		}
		pc.buffer <- buffer[2:]
	}

}

func (pc *AudioServer) readHeader(bh []byte, addr net.Addr) error {
	err := checkHeader(bh)
	if err != nil {
		if err == errMaskInvalid || err == errModeInvalid {
			err := pc.sendStatus(statusNotOK, addr)
			return err
		}
		return err
	}
	err = pc.sendStatus(statusOK, addr)
	if err != nil {
		return err
	}
	return nil
}

func (pc *AudioServer) sendStatus(status byte, addr net.Addr) error {
	msg := []byte{status}
	_, err := pc.pc.WriteTo(msg, addr)
	if err != nil {
		return err
	}
	return nil
}

// Close is to close audio server
func (pc *AudioServer) Close() error {
	if err := pc.pc.Close(); err != nil {
		return err
	}
	return nil
}

func (pc *AudioServer) playAudio() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	// sampleRate := 44100
	sampleRate := 16000
	out := make([]float32, int(16000*1))

	au := &audioOut{buffer: make([]float32, int(16000*1))}
	stream, err := portaudio.OpenDefaultStream(0, 1, float64(sampleRate), len(out), au.processAudio)
	if err != nil {
		return
	}
	defer stream.Close()

	stream.Start()

	pc.wg.Add(1)
	go func(out []float32, au *audioOut) {
		defer pc.wg.Done()
		for v := range pc.buffer {
			r := bytes.NewReader(v)
			err := binary.Read(r, binary.BigEndian, out)
			if err != nil {
				return
			}
			au.buffer = out
			stream.Write()
		}
	}(out, au)
	pc.wg.Wait()
	stream.Stop()

}

type audioOut struct {
	buffer []float32
}

func (a *audioOut) processAudio(out []float32) {
	for i := range a.buffer {
		out[i] = .7 * a.buffer[i] * 4 * 2
	}
}

func checkHeader(r []byte) error {
	if len(r) == 0 {
		return errEmptyByte
	}
	if r[0] != 0x81 {
		return errModeInvalid
	}
	if (r[1] & 0x80) != 0x80 {
		return errMaskInvalid
	}
	return nil
}
