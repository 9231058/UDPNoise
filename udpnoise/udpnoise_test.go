/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 02-05-2018
 * |
 * | File Name:     udpnoise/udpnoise_test.go
 * +===============================================
 */

package udpnoise

import (
	"fmt"
	"net"
	"testing"
)

func TestNoLoss(t *testing.T) {
	const message = "Hello"

	// Destination
	laddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	ln, err := net.ListenUDP("udp", laddr)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := ln.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	t.Logf("Destination on %s", ln.LocalAddr().String())

	// Noise Generator
	us, err := New(0, ln.LocalAddr().String())
	if err != nil {
		t.Fatal(err)
	}

	go us.Run()

	// Source
	raddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", us.Port))
	if err != nil {
		t.Fatal(err)
	}

	ci, err := net.DialUDP("udp", &net.UDPAddr{}, raddr)
	if err != nil {
		t.Fatal(err)
	}

	// Send one packet
	if _, err := ci.Write([]byte(message)); err != nil {
		t.Fatal(err)
	}

	// Receive one packet
	b := make([]byte, len(message))
	if _, err := ln.Read(b); err != nil {
		t.Fatal(err)
	}

	if string(b) != "Hello" {
		t.Fatalf("Send message and received message are not equal: %s != %s", "Hello", b)
	}

	if err := us.Close(); err != nil {
		t.Fatal(err)
	}
}
