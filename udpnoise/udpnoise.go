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
	"fmt"
	"log"
	"math/rand"
	"net"
)

// UDPNoise represents infomation for udp noise proxy instance
type UDPNoise struct {
	Port int

	Loss int

	Destination *net.UDPAddr
	Source      *net.UDPAddr

	ln    *net.UDPConn
	close chan struct{}
}

// New creates new udp noise proxy with given destination and loss probability
func New(loss int, destination string) (*UDPNoise, error) {
	if loss > 100 || loss < 0 {
		return nil, fmt.Errorf("Invalid loss probability: %d is not in [0, 100]", loss)
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

// Run Listen and Forward UDP packets with given loss rate
func (u *UDPNoise) Run() {
	type readUDPData struct {
		data []byte
		from *net.UDPAddr
		err  error
	}

	readUDPChan := make(chan readUDPData)

	go func() {
		for {
			b := make([]byte, 2048)

			n, addr, err := u.ln.ReadFromUDP(b)
			b = b[:n]
			if err != nil {
				readUDPChan <- readUDPData{
					data: nil,
					from: addr,
					err:  err,
				}
			}

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
			readUDPChan <- readUDPData{
				data: b,
				from: addr,
				err:  nil,
			}
		}
	}()

	for {
		// Let's stop the loop
		select {
		case <-u.close:
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

// Close openning udp socket
func (u *UDPNoise) Close() error {
	u.close <- struct{}{}

	return u.ln.Close()
}
