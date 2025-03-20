package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ranjdotdev/httpfromtcp/internal/request"
)

type Server struct {
    listenAddr string
    ln         net.Listener
    quitch     chan struct{}
}

func NewServer(listenAddr string) *Server {
    return &Server{
        listenAddr: listenAddr,
        quitch:     make(chan struct{}),
    }
}

func (s *Server) Start() error {
    ln, err := net.Listen("tcp", s.listenAddr)
    if err != nil {
        return err
    }
    defer ln.Close()
    s.ln = ln

    go s.acceptLoop()

    <-s.quitch

    return nil
}

func (s *Server) acceptLoop() {
    for {
        conn, err := s.ln.Accept()
        if err != nil {
            fmt.Println("accept error: ", err)
            continue
        }
        fmt.Println("new connection: ", conn.RemoteAddr())
        go s.readLoop(conn)
    }
}

func (s *Server) readLoop(conn net.Conn) {
    defer conn.Close()
    req, err := request.RequestFromReader(conn)
    if err != nil {
        fmt.Println("parse error: ", err)
        return
    }
    fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n",
        req.RequestLine.Method,
        req.RequestLine.RequestTarget,
        req.RequestLine.HttpVersion)
    fmt.Println("connection closed: ", conn.RemoteAddr())
}

func main() {
    server := NewServer(":42069")
    log.Fatal(server.Start())
}