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

	ln *net.UDPConn
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

		ln: ln,
	}, nil
}

// Run Listen and Forward UDP packets with given loss rate
func (u *UDPNoise) Run() {
	for {
		b := make([]byte, 1024)

		_, addr, err := u.ln.ReadFromUDP(b)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(addr)

		if rand.Intn(100) < (100 - u.Loss) {
			_, err := u.ln.WriteToUDP(b, u.Destination)
			if err != nil {
				log.Fatal(err)
			}
		}

	}
}

// Close openning udp socket
func (u *UDPNoise) Close() error {
	return u.ln.Close()
}
