package main

import "net"

type User struct {
	id      int
	name    string
	cnn     net.Conn
	groups  map[int]string
	friends []int
}
