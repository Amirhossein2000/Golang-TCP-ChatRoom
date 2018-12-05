package main

import "log"

func errcheck(err error) {
	if err != nil {
		log.Println(err)
	}
}
