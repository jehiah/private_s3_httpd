package main

import (
	"fmt"
	"net/http"
)

type Proxy struct {
	
}

func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "OK")
}