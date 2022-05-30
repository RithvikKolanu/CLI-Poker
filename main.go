package main

import (
	"log"
	"net"
)

func main() {
	s := newServer()

	go s.run()

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Printf("Unable to start server: ", err.Error())
	}

	defer listener.Close()
	log.Println("Started server on :8888")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("unable to accept connection: ", err.Error())
			continue
		}

		go s.newClient(conn)
	}
}
