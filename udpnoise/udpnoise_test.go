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
	"time"
)

func TestMain(t *testing.T) {
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
	t.Logf("Noise Generator on %s", fmt.Sprintf("127.0.0.1:%d", us.Port))

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

	t.Logf("Source on %s", ci.LocalAddr().String())

	t.Run("TestNoLoss", func(t *testing.T) {
		const request = "Hello"
		const response = "My name is 18.20"

		us.Loss = 0

		for i := 0; i < 2; i++ {
			// Send request (source)
			if _, err := ci.Write([]byte(request)); err != nil {
				t.Fatal(err)
			}

			// Receive request (destination)
			bReq := make([]byte, len(request))
			_, addr, err := ln.ReadFromUDP(bReq)
			if err != nil {
				t.Fatal(err)
			}

			if string(bReq) != request {
				t.Fatalf("Send request and received request are not equal: %s != %s", request, bReq)
			}

			// Send response
			if _, err := ln.WriteToUDP([]byte(response), addr); err != nil {
				t.Fatal(err)
			}

			// Receive response (destination)
			bRes := make([]byte, len(response))
			if _, err := ci.Read(bRes); err != nil {
				t.Fatal(err)
			}

			if string(bRes) != response {
				t.Fatalf("Send response and received response are not equal: %s != %s", response, bRes)
			}

		}
	})
	t.Run("TestAllLoss", func(t *testing.T) {
		const request = "GetLoss"

		us.Loss = 100

		// Send request (source)
		if _, err := ci.Write([]byte(request)); err != nil {
			t.Fatal(err)
		}

		// Receive request (destination)
		if err := ln.SetReadDeadline(time.Unix(1, 0)); err != nil {
			t.Fatal(err)
		}
		bReq := make([]byte, len(request))
		_, _, err := ln.ReadFromUDP(bReq)
		if err == nil {
			t.Fatal("This packet must be loss")
		}

	})

	if err := us.Close(); err != nil {
		t.Fatal(err)
	}
}
