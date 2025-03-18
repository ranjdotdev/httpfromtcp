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
  lines := getLinesChannel(data)
  for line := range lines {
    fmt.Printf("read: %s\n", line)
  }
}


func getLinesChannel(d io.ReadCloser) <-chan string {
  lines := make(chan string)
  go func(){
    defer close(lines)

    buffer := make([]byte, 8)
    var currentLine string;
  
    for {
      n, err := d.Read(buffer)
      if err == io.EOF {
        break
      }
      parts := strings.Split(string(buffer[:n]), "\n")
      for i, part := range parts {
        if i != len(parts)-1 {
        currentLine += part
        lines <- currentLine
            currentLine = ""
        } else {
          currentLine += part
        }
      }
    }
    if currentLine != "" {
      lines <- currentLine
    }
  }()
return lines;
}