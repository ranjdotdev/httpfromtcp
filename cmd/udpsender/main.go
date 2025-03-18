package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
  addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
  if err != nil {
      log.Println("Error resolving address:", err)
      return
  }
  conn, err := net.DialUDP("udp", nil, addr)
  if err != nil {
      log.Println("Error creating connection:", err)
      return
  }
  defer conn.Close()
  
  reader := bufio.NewReader(os.Stdin)
  for {
    fmt.Print("> ")

    line, err := reader.ReadString('\n')
    if err != nil {
        log.Println("Error reading input:", err)
        continue
    }

    _, err = conn.Write([]byte(line))
    if err != nil {
        log.Println("Error sending message:", err)
        continue
    }
}
}