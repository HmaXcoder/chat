package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	var adress string
	fmt.Print("Digite o endereço do server: ")
	fmt.Scan(&adress)
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Receber mensagens
	go func() {
		reader := bufio.NewReader(conn)
		for {
			msg, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Desconectado do servidor.")
				os.Exit(0)
			}
			fmt.Print(msg)
		}
	}()

	// Enviar mensagens
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		_, err := fmt.Fprintln(conn, text)
		if err != nil {
			fmt.Println("Erro ao enviar mensagem:", err)
			break
		}
	}
}
