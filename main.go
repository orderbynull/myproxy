package main

import (
	"log"
	"net"
	"io"
)

const MYSQL = "127.0.0.1:3306"
const PROXY = "127.0.0.1:3305"

const COM_QUERY = 3
const COM_STMT_PREPARE = 22

func appToMysql (app net.Conn, mysql net.Conn){
	io.Copy(mysql, app)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	mysql, err := net.Dial("tcp", MYSQL)
	if err != nil{
		log.Fatalf("%s: %s", "ERROR", err.Error())
		return
	}

	go io.Copy(conn, mysql)
	appToMysql(conn, mysql)
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
