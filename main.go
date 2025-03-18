package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type Message struct {
  from    string
  payload string
}

type Server struct{
  listenAddr string
  ln         net.Listener
  quitch     chan struct{}
  msgch      chan Message
}

func NewServer (listenAddr string) *Server {
  return &Server{
    listenAddr: listenAddr,
    quitch: make(chan struct{}),
    msgch:  make(chan Message, 8),
  }
}

func (s* Server) Start() error {
  ln, err := net.Listen("tcp", s.listenAddr)
  if err != nil {return err}
  defer ln.Close();
  s.ln = ln;

  go s.acceptLoop()

  <-s.quitch
  close(s.msgch)

  return nil
}

func (s *Server) acceptLoop () {
  for {
    conn, err := s.ln.Accept()
    if err != nil {
      fmt.Println("accept error: ", err)
      continue;
    }
    fmt.Println("new connection: ", conn.RemoteAddr())
    go s.readLoop(conn)
  }
}

func (s *Server) readLoop(conn net.Conn) {
  defer conn.Close();
    lines := getLinesChannel(conn)
    for line := range lines {
        s.msgch <- Message{
          from: conn.RemoteAddr().String(),
          payload: line,
        }
    }
  fmt.Println("connection closed: ", conn.RemoteAddr())
}


func main() {
  server := NewServer(":42069")
  go func(){
    for msg := range server.msgch {
      fmt.Printf("%s\n", msg.payload)
    }
    }()
  log.Fatal(server.Start())
}


func getLinesChannel(d io.ReadCloser) <-chan string {
  lines := make(chan string)
  go func(){
    defer close(lines)

    buf := make([]byte, 8)
    var currentLine string;
  
    for {
      n, err := d.Read(buf)
      if err == io.EOF {
        break
      }
      parts := strings.Split(string(buf[:n]), "\n")
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