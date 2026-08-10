package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gorilla/handlers"
)

func main() {
	showVersion := flag.Bool("version", false, "print version string")
	listen := flag.String("listen", ":8080", "address:port to listen on.")
	bucket := flag.String("bucket", "", "S3 bucket name")
	logRequests := flag.Bool("log-requests", true, "log HTTP requests")
	region := flag.String("region", "us-east-1", "AWS S3 Region")
	s3Endpoint := flag.String("s3-endpoint", "", "alternate http://address for accessing s3 (for configuring with minio.io)")
	flag.Parse()

	if *showVersion {
		fmt.Printf("private_s3_httpd v%s (built w/%s)\n", VERSION, runtime.Version())
		return
	}

	if *bucket == "" {
		log.Fatalf("bucket name required")
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(*region))
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	var svc *s3.Client
	if *s3Endpoint != "" {
		endpoint := *s3Endpoint
		if !strings.Contains(endpoint, "://") {
			endpoint = "http://" + endpoint
		}
		log.Printf("Using alternate S3 Endpoint with UsePathStyle:true %q", endpoint)
		svc = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = true
		})
	} else {
		svc = s3.NewFromConfig(cfg)
	}

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
