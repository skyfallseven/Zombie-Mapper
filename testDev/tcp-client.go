package main

import (
  "net"
  "fmt"
  "bufio"
  "reflect"
  //"os"
)

func main() {

  // connect to this socket
  conn, _ := net.Dial("tcp", "daisy.student.rit.edu:8081")
  for { 
    // wait for command from server
    message, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
      fmt.Println("Connection to server lost.")
      break
    } else {
      fmt.Print(reflect.TypeOf(message))
    }
        

    /*

    This code has been commented out for dev purposes

     // read in input from stdin
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Text to send: ")
    text, _ := reader.ReadString('\n')

    // send to socket
    fmt.Fprintf(conn, text + "\n")
    */
  }
}