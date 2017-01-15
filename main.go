package main

import (
	"io"
	"log"
	"net"
	dbms "github.com/orderbynull/myproxy/mysql"
	"fmt"
)

const MYSQL = "127.0.0.1:3306"
const PROXY = "127.0.0.1:4040"

func appToMysql(app net.Conn, mysql net.Conn) {
	for{
		pkt, err := dbms.ProxyPacket(app, mysql)
		if err != nil{
			break
		}

		if query, err := dbms.GetQueryString(pkt); err == nil{
			fmt.Printf("> %s \n\n", query)
		}
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	mysql, err := net.Dial("tcp", MYSQL)
	if err != nil {
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
