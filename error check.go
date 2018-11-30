package Golang_TCP_ChatRoom

import "log"

func errcheck(err error) {
	if err != nil {
		log.Println(err)
	}
}
