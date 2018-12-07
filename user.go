package main

import (
	"bufio"
	"net"
)

type User struct {
	id          int
	name        string
	cnn         net.Conn
	chatscanner *bufio.Scanner
	groups      map[int]string
	friends     []int
	chatrequest ChatRequest
}

type ChatRequest struct {
	hostid      int
	guestid     int
	guestaccept bool
}
