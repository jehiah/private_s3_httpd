package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/handlers"
)

func main() {
	showVersion := flag.Bool("version", false, "print version string")
	listen := flag.String("listen", ":8080", "address:port to listen on.")
	bucket := flag.String("bucket", "", "S3 bucket name")
	logRequests := flag.Bool("log-requests", true, "log HTTP requests")
	basicAuthEnabled := flag.Bool("basic-auth", false, "Enforce Basic Auth")
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

	basicAuthUsername, basicAuthUsernameIsSet := os.LookupEnv("PRIVATE_S3_HTTPD_BASIC_AUTH_USERNAME")
	basicAuthPassword, basicAuthPasswordIsSet := os.LookupEnv("PRIVATE_S3_HTTPD_BASIC_AUTH_PASSWORD")

	if *basicAuthEnabled {
		if !basicAuthUsernameIsSet || !basicAuthPasswordIsSet {
			log.Fatalf("Error: Please configure Basic Auth credentials.")
		}
	}

	var svc *s3.S3
	if *s3Endpoint != "" {
		log.Printf("Using alternate S3 Endpoint diwht DisableSSL:true, S3ForcePathStyle:true %q", *s3Endpoint)
		svc = s3.New(session.New(), &aws.Config{
			Region:           region,
			Endpoint:         s3Endpoint,
			DisableSSL:       aws.Bool(true),
			S3ForcePathStyle: aws.Bool(true),
		})
	} else {
		svc = s3.New(session.New(), &aws.Config{
			Region: region,
		})
	}

	var h http.Handler
	h = &Proxy{
		Bucket:            *bucket,
		Svc:               svc,
		basicAuthEnabled:  *basicAuthEnabled,
		basicAuthUsername: basicAuthUsername,
		basicAuthPassword: basicAuthPassword,
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
