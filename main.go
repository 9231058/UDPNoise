package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/elahe-dastan/UDPNoise/udpnoise"
)

func main() {
	fmt.Println("Elahe Dastan <elahe.dstn@gmail.com>")
	fmt.Println("Parham Alvani <parham.alvani@gmail.com>")

	var loss = flag.Int("loss", 0, "packet loss ratio")

	var destination = flag.String("destination", "127.0.0.1:8080", "udp packets destination")

	flag.Parse()

	u, err := udpnoise.New(*loss, *destination)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("UDP Noise on :%d <---> %s with %d loss\n", u.Port, u.Destination, u.Loss)

	u.Run()
}
