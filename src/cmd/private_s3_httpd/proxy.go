package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Proxy struct {
	Bucket            string
	Svc               *s3.S3
	basicAuthEnabled  bool
	basicAuthUsername string
	basicAuthPassword string
}

func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	if p.basicAuthEnabled {

		requestbasicAuthUser, requestbasicAuthPassword, _ := req.BasicAuth()

		if p.basicAuthUsername != requestbasicAuthUser || p.basicAuthPassword != requestbasicAuthPassword {
			rw.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(rw, "Unauthorized.", http.StatusUnauthorized)
			return
		}
	}

	key := req.URL.Path
	if strings.HasSuffix(key, "/") {
		key = key + "index.html"
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(p.Bucket),
		Key:    aws.String(key),
	}
	if v := req.Header.Get("If-None-Match"); v != "" {
		input.IfNoneMatch = aws.String(v)
	}

	// awsReq, resp := p.Svc.GetObjectRequest(input)
	// log.Printf("request: %#v", awsReq)
	// err := awsReq.Send()
	// log.Printf("response: %#v", )

	var is304 bool
	resp, err := p.Svc.GetObject(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case s3.ErrCodeNoSuchKey:
			http.Error(rw, "Page Not Found", 404)
			return
		case "NotModified":
			is304 = true
			// continue so other headers get set appropriately
		default:
			log.Printf("Error: %v %v %v", awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
			http.Error(rw, "Internal Error", 500)
			return
		}
	} else if err != nil {
		log.Printf("not aws error %v %s", err, err)
		http.Error(rw, "Internal Error", 500)
		return
	}

	var contentType string
	if resp.ContentType != nil {
		contentType = *resp.ContentType
	}

	if contentType == "" {
		ext := path.Ext(req.URL.Path)
		contentType = mime.TypeByExtension(ext)
	}

	if resp.ETag != nil && *resp.ETag != "" {
		rw.Header().Set("Etag", *resp.ETag)
	}

	if contentType != "" {
		rw.Header().Set("Content-Type", contentType)
	}
	if resp.ContentLength != nil && *resp.ContentLength > 0 {
		rw.Header().Set("Content-Length", fmt.Sprintf("%d", *resp.ContentLength))
	}

	if is304 {
		rw.WriteHeader(304)
	} else {
		io.Copy(rw, resp.Body)
		resp.Body.Close()
	}
}

// resp, err := svc.ListObjects(&s3.ListObjectsInput{
// 	Bucket:  aws.String(settings.GetString("s3_bucket")),
// 	Prefix:  aws.String("data/"),
// 	MaxKeys: aws.Long(1000),
// })
// if awsErr, ok := err.(awserr.Error); ok {
// 	// A service error occurred.
// 	log.Fatalf("Error: %v %v", awsErr.Code, awsErr.Message)
// } else if err != nil {
// 	// A non-service error occurred.
// 	log.Fatalf("%v", err)
// }
