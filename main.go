package Golang_TCP_ChatRoom

import (
	"log"
	"math/rand"
	"net"
	"time"
)

const port = "9000"

func main() {

	rand.Seed(time.Now().Unix())

	ln, err := net.Listen("tcp", ":"+port)

	log.Println("lisening on",port)

	errcheck(err)
	defer ln.Close()
	for {
		cnn, err := ln.Accept()
		errcheck(err)
		go handle(cnn)
	}

	log.Println("Connection Ended")

}
