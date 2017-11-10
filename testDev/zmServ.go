/*
Testing C2 Server
This is the command
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
Name: 		HandleFunc
Descript:	Handles incoming command
Param:		Receives open connection in a ReadWriter interface
*/
type HandleFunc func(*bufio.ReadWriter)

/*
Name:		Endpoint
Descript:	Provides an "endpoint" to other processes for data
*/ 
type Endpoint struct {
	listener net.listener 	//New listeners for multiple conns
	handler map[string]HandleFunc //each command has a diff func to call

	//Maps are not threadsafe, we need a mutex to control access
	m sync.RWMutex
}

/*
Name: 		NewEndpoint
Returns:	Endpoint object for a new function/set of processes
*/
func NewEndpoint() *Endpoint {
	return &Endpoint {handler: map[string]HandleFunc{},}
}

/*
Name: 	AddHandleFunc
Desc:	Adds new function for handling data
Note:	I imagine this helps with scalability of commands
*/
func (e *Endpoint) AddHandleFunc(name string, f HandleFunc) {
	e.m.Lock()
	e.hanler[name] = f
	e.m.Unlock()
}

/*
Name: 	Listen
Desc:	Listen on the endpoint port on all interfaces
Pre:	One handler function exists via AddHandleFunc()
*/
func (e *Endpoint) Listen() error {
	var	err error
	e.listener, err = net.Listen("tcp", Port)
	if err != nil {
		return errors.Wrap(err, "ListenerFailure\n")
	}
	log.Prinln("Listening on", e.Listener.Addr().String())

	for {
		log.Println("Awaiting clients to connect")
		conn, err := e.Listener.Accept()
		if err != nil {
			log.Println("ConnectionRequest Failure:", err)
			continue
		}
		log.Println("zmSession starting...")
		go e.handleMessages(conn)
	}
}
