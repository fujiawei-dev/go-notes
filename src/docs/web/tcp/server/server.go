package main

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	_ "net/http/pprof"
)

func main() {

	ln, err := net.Listen("tcp", ":8972")
	if err != nil {
		panic(err)
	}

	var connections []net.Conn
	defer func() {
		for _, conn := range connections {
			conn.Close()
		}
	}()

	for {
		conn, e := ln.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				log.Printf("accept temp err: %v", ne)
				continue
			}

			log.Printf("accept err: %v", e)
			return
		}

		go handleConn(conn)
		connections = append(connections, conn)
		log.Printf("total number of connections: %v", len(connections))
	}
}

func handleConn(conn net.Conn) {
	io.Copy(ioutil.Discard, conn)
	//io.Copy(os.Stdout, conn)
}
