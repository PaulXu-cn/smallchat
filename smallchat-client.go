package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"golang.org/x/term"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <host> <port>\n", os.Args[0])
		os.Exit(1)
	}

	conn, err := net.Dial("tcp", os.Args[1]+":"+os.Args[2])
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()

	var tty = term.NewTerminal(conn, "> ")

	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	// get stdin in
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		tty.Write([]byte(scanner.Text() + "\n"))
	}

	return
}
