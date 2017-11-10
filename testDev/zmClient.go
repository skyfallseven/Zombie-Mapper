/*
Testing C2 Server
This is the bot running the scanner
*/

package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"encoding/gob"
	"flag"
)

const (Port = ":61000")

/*
Name: 		Open
Descript:	Connects to a TCP address
Params: 	addr: the tcp address to connect to
Returns:	TCP connection w/timeout wrapped in a ReadWriter
*/
func Open(addr string) (*bufio.ReadWriter, error) {

	log.Println("Connecting to " + addr)
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return nil, errors.Wrap(err, "Connection failed.")
	}

	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}