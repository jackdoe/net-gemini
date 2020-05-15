// gemini server implementation, check out gemini://gemini.circumlunar.space/
// inspired by https://tildegit.org/solderpunk/molly-brown
//
// Example:
//
//  package main
//
//  import (
//  	"log"
//
//  	gemini "github.com/jackdoe/net-gemini"
//  )
//
//  func main() {
//	gemini.HandleFunc("/example", func(w *gemini.Response, r *gemini.Request) {
//		w.SetStatus(gemini.StatusSuccess, "text/gemini")
//		w.Write([]byte("HELLO: " + r.URL.Path + "\n"))
//	})
//
//	gemini.Handle("/", gemini.FileServer(*root))
//  	log.Fatal(gemini.ListenAndServeTLS("localhost:1965", "localhost.crt", "localhost.key"))
//  }
//
// You can also checkout cmd/main.go as an example.
//
// Make sure to generate your cert for localhost:
//
//  openssl req \
//          -x509 \
//          -out localhost.crt \
//          -keyout localhost.key \
//          -newkey rsa:2048 \
//          -nodes \
//          -sha256 \
//          -subj '/CN=localhost' \
//          -extensions EXT \
//          -config <( printf "[dn]\nCN=localhost\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:localhost\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")
//
package gemini

import (
	"os"
	"strings"
)

var srv = &Server{}

type HandlerFunc func(*Response, *Request)

func (f HandlerFunc) ServeGemini(w *Response, r *Request) {
	f(w, r)
}

type handledPath struct {
	handler Handler
	p       string
}

type basicHandler struct {
	handlers []handledPath
}

func (b *basicHandler) ServeGemini(w *Response, r *Request) {
	u := r.URL.Path
	for _, h := range b.handlers {
		if strings.HasPrefix(u, h.p) {
			h.handler.ServeGemini(w, r)
			return
		}
	}
	w.SetStatus(StatusNotFound, u+" Not Found!")
}

var basic = &basicHandler{}

func Handle(p string, h Handler) {
	basic.handlers = append(basic.handlers, handledPath{handler: h, p: p})
}

func HandleFunc(p string, f HandlerFunc) {
	basic.handlers = append(basic.handlers, handledPath{handler: f, p: p})
}

func ListenAndServeTLS(addr string, certFile, keyFile string) error {
	s := Server{Addr: addr, Handler: basic, Log: os.Stdout}
	return s.ListenAndServeTLS(certFile, keyFile)
}
