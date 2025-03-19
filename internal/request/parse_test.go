package request

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
    data            string
    numBytesPerRead int
    pos             int
}

func (cr *chunkReader) Read(p []byte) (n int, err error) {
    if cr.pos >= len(cr.data) {
        return 0, io.EOF
    }
    endIndex := cr.pos + cr.numBytesPerRead
    if endIndex > len(cr.data) {
        endIndex = len(cr.data)
    }
    n = copy(p, cr.data[cr.pos:endIndex])
    cr.pos += n
    if n > cr.numBytesPerRead {
        n = cr.numBytesPerRead
        cr.pos -= n - cr.numBytesPerRead
    }
    return n, nil
}

func TestParseGoodGetRequestChunkSize3(t *testing.T) {
    req := &Request{State: StateInitialized}
    reader := &chunkReader{
        data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\n",
        numBytesPerRead: 3,
    }
    buf := make([]byte, 8)

    for {
        n, err := reader.Read(buf)
        if n > 0 {
            consumed, parseErr := req.parse(buf[:n])
            require.NoError(t, parseErr)
            if req.State == StateDone {
                assert.Equal(t, "GET", req.RequestLine.Method)
                assert.Equal(t, "/", req.RequestLine.RequestTarget)
                assert.Equal(t, "1.1", req.RequestLine.HttpVersion)
                assert.Equal(t, 14, consumed)
                break
            }
        }
        if err == io.EOF {
            t.Fatal("should have parsed before EOF")
        }
        if err != nil {
            t.Fatal(err)
        }
    }
}

func TestParseGoodGetRequestWithPathChunkSize1(t *testing.T) {
    req := &Request{State: StateInitialized}
    reader := &chunkReader{
        data:            "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\n",
        numBytesPerRead: 1,
    }
    buf := make([]byte, 8)

    for {
        n, err := reader.Read(buf)
        if n > 0 {
            consumed, parseErr := req.parse(buf[:n])
            require.NoError(t, parseErr)
            if req.State == StateDone {
                assert.Equal(t, "GET", req.RequestLine.Method)
                assert.Equal(t, "/coffee", req.RequestLine.RequestTarget)
                assert.Equal(t, "1.1", req.RequestLine.HttpVersion)
                assert.Equal(t, 20, consumed)
                break
            }
        }
        if err == io.EOF {
            t.Fatal("should have parsed before EOF")
        }
        if err != nil {
            t.Fatal(err)
        }
    }
}

func TestParseInvalidPartsChunkSize3(t *testing.T) {
    req := &Request{State: StateInitialized}
    reader := &chunkReader{
        data:            "/coffee HTTP/1.1\r\nHost: localhost:42069\r\n",
        numBytesPerRead: 3,
    }
    buf := make([]byte, 8)

    for {
        n, err := reader.Read(buf)
        if n > 0 {
            _, parseErr := req.parse(buf[:n])
            if parseErr != nil {
                require.Error(t, parseErr)
                return
            }
            if req.State == StateDone {
                t.Fatal("should have errored on invalid parts")
            }
        }
        if err == io.EOF {
            t.Fatal("should have errored before EOF")
        }
        if err != nil {
            t.Fatal(err)
        }
    }
}