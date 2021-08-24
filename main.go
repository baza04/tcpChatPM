package main

import (
	"flag"
	"log"
	"net"
)

func main() {
	port := flag.String("port", "8080", "custom port to tcp chat")
	flag.Parse()

	s := newServer()
	go s.run()

	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}
	defer listener.Close()
	log.Println("listen on: ", *port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("unable to accept connection: %s", err.Error())
			continue
		}

		go s.newClient(conn)
	}
}
