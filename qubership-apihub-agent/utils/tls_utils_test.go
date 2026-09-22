package utils

import (
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestCloseAfterRequestTimesOutClosesTCP(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	serverClosed := make(chan struct{})
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		buf := make([]byte, 4096)
		_, _ = conn.Read(buf)
		_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
		_, err = conn.Read(buf)
		if err != nil {
			close(serverClosed)
		}
	}()

	client, err := CreateSecureHTTPClientCloseAfterRequest(200 * time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Get("http://" + ln.Addr().String() + "/")
	if err == nil {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
		t.Fatal("expected client timeout, got success")
	}
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}

	select {
	case <-serverClosed:
	case <-time.After(3 * time.Second):
		t.Fatal("server connection stayed open after client timeout")
	}
}
