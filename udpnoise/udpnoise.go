/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 01-05-2018
 * |
 * | File Name:     udpnoise.go
 * +===============================================
 */

package udpnoise

import (
	"errors"
	"log"
	"math/rand"
	"net"
)

var (
	ErrInvalidLoss = errors.New("invalid loss probability, it must be in [0, 100]")
)

// Message that is read from udp socket.
type Message struct {
	data []byte
	from *net.UDPAddr
	err  error
}

// UDPNoise represents information for udp noise proxy instance.
// peer -> source --/-- destination --> peer.
// peer <- source ----- destination <-- peer.
type UDPNoise struct {
	// Port is the source port and one of the peers must connect to it.
	Port int

	// Loss is the loss ratio for ongoing packetes.
	Loss int

	// Destination address
	Destination *net.UDPAddr
	// Source address
	Source *net.UDPAddr

	ln *net.UDPConn

	close chan struct{}
}

// New creates new udp noise proxy with given destination and loss probability.
func New(loss int, destination string) (*UDPNoise, error) {
	if loss > 100 || loss < 0 {
		return nil, ErrInvalidLoss
	}

	addr, err := net.ResolveUDPAddr("udp", destination)
	if err != nil {
		return nil, err
	}

	ln, err := net.ListenUDP("udp", &net.UDPAddr{})
	if err != nil {
		return nil, err
	}

	return &UDPNoise{
		Port: ln.LocalAddr().(*net.UDPAddr).Port,

		Loss: loss,

		Destination: addr,
		Source:      nil,

		ln:    ln,
		close: make(chan struct{}),
	}, nil
}

func (u *UDPNoise) reader() <-chan Message {
	readUDPChan := make(chan Message)

	go func() {
		for {
			b := make([]byte, 2048)

			n, addr, err := u.ln.ReadFromUDP(b)
			if err != nil {
				readUDPChan <- Message{
					data: nil,
					from: addr,
					err:  err,
				}
			}

			b = b[:n]

			// store source address
			if addr.String() != u.Destination.String() {
				if u.Source != nil {
					if u.Source.String() != addr.String() {
						u.Source = addr
					}
				} else {
					u.Source = addr
				}
			}

			log.Printf("[udpnoise] Packet from %s", addr)
			readUDPChan <- Message{
				data: b,
				from: addr,
				err:  nil,
			}
		}
	}()

	return readUDPChan
}

// Run Listen and Forward UDP packets with given loss rate.
func (u *UDPNoise) Run() {
	readUDPChan := u.reader()

	for {
		select {
		case <-u.close: // close the read loop
			return
		case d := <-readUDPChan:
			if d.err != nil {
				log.Fatalf("[udpnoise] Read from UDP: %s", d.err)
			}

			if rand.Intn(100) < (100 - u.Loss) {
				if d.from.String() != u.Destination.String() {
					if _, err := u.ln.WriteToUDP(d.data, u.Destination); err != nil {
						log.Fatalf("[udpnoise] Write to UDP (%s): %s", u.Destination, err)
					}

					log.Printf("[udpnoise] Packet sends to %s with loss rate %d", u.Destination, u.Loss)
				} else {
					if _, err := u.ln.WriteToUDP(d.data, u.Source); err != nil {
						log.Fatalf("[udpnoise] Write to UDP (%s): %s", u.Source, err)
					}

					log.Printf("[udpnoise] Packet sends to %s with loss rate %d", u.Source, u.Loss)
				}
			}
		}
	}
}

// Close openning udp socket.
func (u *UDPNoise) Close() error {
	close(u.close)
	return u.ln.Close()
}
