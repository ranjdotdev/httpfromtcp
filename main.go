package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
  data, err := os.Open("messages.txt")
  if err != nil { log.Fatal("error: ", err) }
  defer data.Close()
  
  buffer := make([]byte, 8)
  var currentLine string;
  
  for {
// 1st
    n, err := data.Read(buffer)
    if err == io.EOF {
      break
    }
// 2nd
parts := strings.Split(string(buffer[:n]), "\n")
// 3rd
for i, part := range parts {
  if i != len(parts)-1 {
    currentLine += part
    fmt.Printf("read: %s\n", string(currentLine))
    currentLine = ""
    } else {
      currentLine += part
    }
  }
}
  if currentLine != "" {
    fmt.Printf("read: %s\n", string(currentLine))
  }
}
