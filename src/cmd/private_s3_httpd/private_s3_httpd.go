package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/handlers"
)

func main() {
	listen := flag.String("listen", ":8080", "address:port to listen on.")
	bucket := flag.String("bucket", "", "S3 bucket name")
	logRequests := flag.Bool("log-requests", true, "log HTTP requests")
	region := flag.String("region", "us-east-1", "AWS S3 Region")
	flag.Parse()

	if *bucket == "" {
		log.Fatalf("bucket name required")
	}

	svc := s3.New(session.New(), &aws.Config{
		Region: region,
	})

	var h http.Handler
	h = &Proxy{
		Bucket: *bucket,
		Svc:    svc,
	}
	if *logRequests {
		h = handlers.CombinedLoggingHandler(os.Stdout, h)
	}

	s := &http.Server{
		Addr:           *listen,
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("listening on %s", *listen)
	log.Fatal(s.ListenAndServe())

}
