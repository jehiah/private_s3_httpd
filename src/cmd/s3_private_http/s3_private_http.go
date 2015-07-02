package main

import (
	"net/http"
	"time"
	"log"
	"flag"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	listen := flag.String("listen", ":8080", "address:port to listen on.")
	bucket := flag.String("bucket", "", "s3 bucket name")
	flag.Parse()
	
	if *bucket == "" {
		log.Fatalf("bucket name required")
	}


	svc := s3.New(&aws.Config{
		// Credentials: creds,
		Region: "us-east-1",
	})
	
	p := &Proxy{
		Bucket: *bucket,
		Svc: svc,
	}
	
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