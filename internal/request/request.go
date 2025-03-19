package request

import (
	"errors"
	"io"
	"strings"
	"unicode"
)

type Request struct {
    RequestLine RequestLine
    State       uint8
    Buffer      string
}

type RequestLine struct {
    HttpVersion   string
    RequestTarget string
    Method        string
}

const (
    StateInitialized = 0
    StateDone        = 1
)

func (r *Request) parse(data []byte) (int, error) {
    r.Buffer = r.Buffer + string(data)
    if strings.Contains(r.Buffer, "\r\n") {
        line := strings.Split(r.Buffer, "\r\n")[0]
        parts := strings.Split(line, " ")
        if len(parts) != 3 {
            return 0, errors.New("request line must have 3 parts")
        }
        method := parts[0]
        target := parts[1]
        version := parts[2]

        for _, char := range method {
            if !unicode.IsUpper(char) || !unicode.IsLetter(char) {
                return 0, errors.New("method must be uppercase letters")
            }
        }
        if version != "HTTP/1.1" {
            return 0, errors.New("only HTTP/1.1 is supported")
        }

        r.RequestLine.Method = method
        r.RequestLine.RequestTarget = target
        r.RequestLine.HttpVersion = "1.1"
        r.State = StateDone
        return len(line), nil
    }
    return 0, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
    req := &Request{
        State:  StateInitialized,
        Buffer: "",
    }
    chunk := make([]byte, 8)
    for req.State != StateDone {
        n, err := reader.Read(chunk)
        if n > 0 {
            consumed, err := req.parse(chunk[:n])
            if err != nil {
                return nil, err
            }
            if consumed > 0 {
                req.Buffer = req.Buffer[consumed:]
            }
        }
        if err == io.EOF {
            if req.State != StateDone {
                return nil, errors.New("incomplete request")
            }
            break
        }
        if err != nil {
            return nil, err
        }
    }
    return req, nil
}