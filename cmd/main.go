package main

import (
	"flag"
	"log"

	gemini "github.com/jackdoe/net-gemini"
)

func main() {
	root := flag.String("root", "", "root directory")
	crt := flag.String("crt", "", "path to cert")
	key := flag.String("key", "", "path to cert key")
	bind := flag.String("bind", "localhost:1965", "bind to")
	flag.Parse()

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

	log.Fatal(gemini.ListenAndServeTLS(*bind, *crt, *key))
}
