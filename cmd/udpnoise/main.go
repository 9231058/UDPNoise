/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 02-05-2018
 * |
 * | File Name:     main.go
 * +===============================================
 */

package main

import (
	"fmt"
	"log"

	"github.com/aut-ceit/UDPNoise/udpnoise"
)

func main() {
	fmt.Println("Parham Alvani <parham.alvani@gmail.com>")

	loss := 100
	destination := "127.0.0.1:8080"

	u, err := udpnoise.New(loss, destination)
	if err != nil {
		log.Fatal(err)
	}

	u.Run()
}
