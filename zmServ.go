/*
Zombie Mapper [SERVER] v0.001-DEV
Control the zm sniffer from a remote interface
Currently using HTTP, will switch to HTTPS upon further development
Possibly use custom packet development
*/

package main

import (
	"net"
	"fmt"
	"bufio"
	"os"
)

func main() {
	fmt.Println("Zombie Mapper Control Server")

	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":8081")

	// accept connection on port
}