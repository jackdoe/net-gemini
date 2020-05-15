package gemini

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/url"
	"time"
)

type Handler interface {
	ServeGemini(*Response, *Request)
}

type Server struct {
	Addr         string
	Handler      Handler
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	listener     net.Listener
	Log          io.Writer
}

func (s *Server) logf(format string, args ...interface{}) {
	if s.Log != nil {
		now := fmt.Sprintf("%v ", time.Now().Format(time.ANSIC))
		fmt.Fprintf(s.Log, now+format+"\n", args...)
	}
}

func (s *Server) ListenAndServeTLS(certFile, keyFile string) error {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}
	tlscfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	listener, err := tls.Listen("tcp", s.Addr, tlscfg)
	if err != nil {
		return err
	}
	defer listener.Close()

	s.logf("gemini listening on %s [tls: %v %v]", listener.Addr(), certFile, keyFile)
	return s.serve(listener)
}

func (s *Server) serve(listener net.Listener) error {
	s.listener = listener

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go s.handleGeminiRequest(conn)
	}
}

func (s *Server) handleGeminiRequest(conn net.Conn) {
	readDeadline := time.Time{}
	in := Request{}
	out := Response{conn: conn}
	t0 := time.Now()
	if d := s.ReadTimeout; d != 0 {
		readDeadline = t0.Add(d)
		conn.SetReadDeadline(readDeadline)
	}
	// FIXME: this is read + write, should be only write
	if d := s.WriteTimeout; d != 0 {
		conn.SetWriteDeadline(time.Now().Add(d))
	}

	// FIXME: use something else
	defer func() {
		s.logf("%s -> %s request: %s took %v, status: %v %v", conn.RemoteAddr(), conn.LocalAddr(), in.URL, time.Since(t0), out.statusCode, out.statusText)
	}()

	defer conn.Close()

	reader := bufio.NewReaderSize(conn, 1024)
	request, overflow, err := reader.ReadLine()
	if overflow {
		_ = out.SetStatus(StatusPermanentFailure, "Request too long!")
		return
	} else if err != nil {
		_ = out.SetStatus(StatusTemporaryFailure, "Unknown error reading request! "+err.Error())
		return
	}

	URL, err := url.Parse(string(request))
	if err != nil {
		_ = out.SetStatus(StatusPermanentFailure, "Error parsing URL! "+err.Error())
		return
	}
	if URL.Scheme == "" {
		URL.Scheme = "gemini"
	}

	if URL.Scheme != "gemini" {
		_ = out.SetStatus(StatusPermanentFailure, "No proxying to non-Gemini content!")
		return
	}
	in.URL = URL

	s.Handler.ServeGemini(&out, &in)
}

func (s *Server) Close() {
	s.listener.Close()
}
