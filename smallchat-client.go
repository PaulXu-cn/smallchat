package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

// Terminal settings
type Termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Line   uint8
	Cc     [32]uint8
	Ispeed uint32
	Ospeed uint32
}

const (
	TCGETS = 0x5401
	TCSETS = 0x5402
)

func getTermios(fd int) (*Termios, error) {
	termios := &Termios{}
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), TCGETS, uintptr(unsafe.Pointer(termios)))
	if err != 0 {
		return nil, err
	}
	return termios, nil
}

func setTermios(fd int, termios *Termios) error {
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), TCSETS, uintptr(unsafe.Pointer(termios)))
	if err != 0 {
		return err
	}
	return nil
}

func setRawMode(fd int) error {
	termios, err := getTermios(fd)
	if err != nil {
		return err
	}

	termios.Iflag &^= syscall.BRKINT | syscall.ICRNL | syscall.INPCK | syscall.ISTRIP | syscall.IXON
	termios.Oflag &^= syscall.OPOST
	termios.Cflag |= syscall.CS8
	termios.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG
	termios.Cc[syscall.VMIN] = 1
	termios.Cc[syscall.VTIME] = 0

	return setTermios(fd, termios)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <host> <port>\n", os.Args[0])
		os.Exit(1)
	}

	port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid port number")
		os.Exit(1)
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", os.Args[1], port))
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	if err := setRawMode(int(os.Stdin.Fd())); err != nil {
		fmt.Println("Error setting raw mode:", err)
		os.Exit(1)
	}

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text()
			if _, err := conn.Write([]byte(text + "\n")); err != nil {
				fmt.Println("Error writing to server:", err)
				os.Exit(1)
			}
		}
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Println(strings.TrimRight(text, "\n"))
	}
}
