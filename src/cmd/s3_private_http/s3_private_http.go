package main

import (
	"net/http"
	"time"
	"log"
	"flag"
)

func main() {
	listen := flag.String("listen", ":8080", "address:port to listen on.")
	flag.Parse()
	
	p := &Proxy{}
	
	s := &http.Server{
		Addr:           *listen,
		Handler:        p,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
    log.Printf("listening on %s", *listen)
	log.Fatal(s.ListenAndServe())
	
}