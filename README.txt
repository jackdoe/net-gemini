package gemini // import "github.com/jackdoe/net-gemini"

gemini server implementation, check out gemini://gemini.circumlunar.space/
inspired by https://tildegit.org/solderpunk/molly-brown

Example:

     package main

     import (
     	"log"

     	gemini "github.com/jackdoe/net-gemini"
     )

     func main() {
    	gemini.HandleFunc("/example", func(w *gemini.Response, r *gemini.Request) {
    		if len(r.URL.RawQuery) == 0 {
    			w.SetStatus(gemini.StatusInput, "what is the answer to the ultimate question")
    		} else {
    			w.SetStatus(gemini.StatusSuccess, "text/gemini")
    			answer := r.URL.RawQuery
    			w.Write([]byte("HELLO: " + r.URL.Path + ", yes the answer is: " + answer))
    		}
    	})

    	gemini.Handle("/", gemini.FileServer(*root))
     	log.Fatal(gemini.ListenAndServeTLS("localhost:1965", "localhost.crt", "localhost.key"))
     }

You can also checkout cmd/main.go as an example.

Make sure to generate your cert for localhost:

    openssl req \
            -x509 \
            -out localhost.crt \
            -keyout localhost.key \
            -newkey rsa:2048 \
            -nodes \
            -sha256 \
            -subj '/CN=localhost' \
            -extensions EXT \
            -config <( printf "[dn]\nCN=localhost\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:localhost\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")

CONSTANTS

const (
	StatusInput                              = 10
	StatusSuccess                            = 20
	StatusSuccessEndClientCertificateSession = 21
	StatusRedirectTemporary                  = 30
	StatusRedirectPermanent                  = 31
	StatusTemporaryFailure                   = 40
	StatusServerUnavailable                  = 41
	StatusCGIError                           = 42
	StatusProxyError                         = 43
	StatusSlowDown                           = 44
	StatusPermanentFailure                   = 50
	StatusNotFound                           = 51
	StatusGone                               = 52
	StatusProxyRequestRefused                = 53
	StatusBadRequest                         = 59
	StatusClientCertRequired                 = 60
	StatusTransientCertRequested             = 61
	StatusAuthorisedCertRequired             = 62
	StatusCertNotAccepted                    = 63
	StatusFutureCertRejected                 = 64
	StatusExpiredCertRejected                = 65
)

FUNCTIONS

func Handle(p string, h Handler)
func HandleFunc(p string, f HandlerFunc)
func ListenAndServeTLS(addr string, certFile, keyFile string) error
func ServeFilePath(p string, w *Response, r *Request)

TYPES

type Handler interface {
	ServeGemini(*Response, *Request)
}

func FileServer(root string) Handler

type HandlerFunc func(*Response, *Request)

func (f HandlerFunc) ServeGemini(w *Response, r *Request)

type Request struct {
	URL *url.URL
}

type Response struct {
	// Has unexported fields.
}

func (r *Response) SetStatus(s Status, mime string) error

func (r *Response) Write(b []byte) (int, error)

type Server struct {
	Addr         string
	Handler      Handler
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	Log io.Writer
	// Has unexported fields.
}

func (s *Server) Close()

func (s *Server) ListenAndServeTLS(certFile, keyFile string) error

type Status int

