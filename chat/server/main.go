package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

var clients = make(map[net.Conn]string)
var mutex = sync.Mutex{}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	mutex.Lock()
	clients[conn] = conn.RemoteAddr().String()
	mutex.Unlock()

	fmt.Println("Novo cliente:", conn.RemoteAddr())

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		msg = strings.TrimSpace(msg)

		broadcast(fmt.Sprintf("%s: %s\n", clients[conn], msg), conn)
	}

	mutex.Lock()
	delete(clients, conn)
	mutex.Unlock()
	fmt.Println("Cliente saiu:", conn.RemoteAddr())
}

func broadcast(message string, sender net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	for conn := range clients {
		if conn != sender {
			conn.Write([]byte(message))
		}
	}
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("Servidor rodando na porta 8080...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}
