package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
  data, err := os.Open("messages.txt")
  if err != nil { log.Fatal("error: ", err) }
  defer data.Close()

  buffer := make([]byte, 8)
  for {
    n, err := data.Read(buffer)
    if err == io.EOF {
      break
    }
    fmt.Printf("read: %s\n", string(buffer[:n]))
  }
}