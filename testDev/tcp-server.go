package main

import (
  "net"
  "fmt"
  "bufio"
  //"strings" 
  "os"
)

func main() {

  fmt.Println("Launching server...")

  // listen on all interfaces
  ln, _ := net.Listen("tcp", ":8081")

  // accept connection on port
  conn, _ := ln.Accept()

  // run loop forever (or until ctrl-c)
  for {

    reader := bufio.NewReader(os.Stdin)
    fmt.Print("> ")
    cmd, _ := reader.ReadString('\n')
    fmt.Fprintf(conn, cmd, '\n')


    // will listen for message to process ending in newline (\n)
    //message, _ := bufio.NewReader(conn).ReadString('\n')

    // output message received
    //fmt.Print("Message Received:", string(message))

    // sample process for string received
    // newmessage := strings.ToUpper(message)

    // send new string back to client
    // conn.Write([]byte(newmessage + "\n"))
  }
}