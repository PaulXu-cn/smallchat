package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

// Set the specified socket in non-blocking mode, with no delay flag.
func socketSetNonBlockNoDelay(conn net.Conn) error {
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return fmt.Errorf("not a TCP connection")
	}

	if err := tcpConn.SetNoDelay(true); err != nil {
		return err
	}

	if err := tcpConn.SetReadDeadline(time.Now()); err != nil {
		return err
	}

	return nil
}

// Create a TCP socket listening to 'port' ready to accept connections.
func createTCPServer(port int) (net.Listener, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	return listener, nil
}

// Create a TCP socket and connect it to the specified address.
func TCPConnect(addr string, port int, nonblock bool) (net.Conn, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return nil, err
	}

	if nonblock {
		if err := socketSetNonBlockNoDelay(conn); err != nil {
			conn.Close()
			return nil, err
		}
	}

	return conn, nil
}

// If the listening socket signaled there is a new connection ready to be accepted, we accept it.
func acceptClient(listener net.Listener) (net.Conn, error) {
	conn, err := listener.Accept()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// We also define an allocator that always crashes on out of memory.
func chatMalloc(size int) []byte {
	data := make([]byte, size)
	if data == nil {
		fmt.Println("Out of memory")
		os.Exit(1)
	}

	return data
}

// Also aborting realloc().
func chatRealloc(slice []byte, size int) []byte {
	newSlice := make([]byte, size)
	if newSlice == nil {
		fmt.Println("Out of memory")
		os.Exit(1)
	}

	copy(newSlice, slice)
	return newSlice
}
