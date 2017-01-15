package main

import (
	"log"
	"net"
	"io"
)

const MYSQL = "127.0.0.1:3306"
const PROXY = "127.0.0.1:3305"

func handleConnection(conn net.Conn) {
	defer conn.Close()

	mysql, err := net.Dial("tcp", MYSQL)
	if err != nil{
		log.Fatalf("%s: %s", "ERROR", err.Error())
		return
	}

	go io.Copy(conn, mysql)
	io.Copy(mysql, conn)

}

func main() {
	proxy, err := net.Listen("tcp", PROXY)
	if err != nil {
		log.Fatalf("%s: %s", "ERROR", err.Error())
	}
	defer proxy.Close()

	for {
		conn, err := proxy.Accept()
		if err != nil {
			log.Printf("%s: %s", "ERROR", err.Error())
		}

		go handleConnection(conn)
	}
}
