package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

const (
	MaxClients = 1000
	ServerPort = "7711"
)

type Client struct {
	conn net.Conn
	nick string
}

type ChatState struct {
	serverConn net.Listener
	clients    map[net.Conn]*Client
	mu         sync.Mutex
	numclients int
	maxclient  int
}

var Chat = &ChatState{
	clients: make(map[net.Conn]*Client),
}

func createClient(fd net.Conn) *Client {
	nick := fmt.Sprintf("user:%s", fd.RemoteAddr().String())
	c := &Client{
		conn: fd,
		nick: nick,
	}
	Chat.clients[fd] = c

	Chat.numclients++

	return c
}

func freeClient(c *Client) {
	c.conn.Close()
	delete(Chat.clients, c.conn)
	Chat.numclients--
}

func initChat() {
	var err error
	Chat.serverConn, err = net.Listen("tcp", ":"+ServerPort)
	if err != nil {
		fmt.Println("Error creating server:", err)
		os.Exit(1)
	}
}

func sendMsgToAllClientsBut(excluded net.Conn, msg string) {
	Chat.mu.Lock()
	defer Chat.mu.Unlock()
	for conn, client := range Chat.clients {
		if conn == excluded {
			continue
		}
		fmt.Fprintln(client.conn, msg)
	}
}

func handleClient(c *Client) {
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "/nick ") {
			c.nick = strings.TrimPrefix(text, "/nick ")
		} else {
			msg := fmt.Sprintf("%s> %s", c.nick, text)
			fmt.Println(msg)
			sendMsgToAllClientsBut(c.conn, msg)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from client:", err)
	}
	freeClient(c)
}

func main() {
	initChat()
	fmt.Print("Server started on port " + ServerPort + "\n")

	for {
		conn, err := Chat.serverConn.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		c := createClient(conn)
		fmt.Fprintln(c.conn, "Welcome to Simple Chat! Use /nick <nick> to set your nick.\n")
		fmt.Printf("Connected client  %s\n", c.conn.LocalAddr().String())
		go handleClient(c)
	}
}
