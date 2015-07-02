package main

import (
	"fmt"
	"io"
	"mime"
	"log"
	"net/http"
	"path"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Proxy struct {
	Bucket string
	Svc    *s3.S3
}

func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	key := req.URL.Path
	if key == "/" {
		key = "/index.html"
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
	
	resp, err := p.Svc.GetObject(input)
	if awsErr, ok := err.(awserr.Error); ok {
		if awsErr.Code() == "NoSuchKey" {
			http.Error(rw, "Page Not Found", 404)
		} else {
			// A service error occurred.
			log.Printf("Error: %v %v", awsErr.Code(), awsErr.Message())
			http.Error(rw, "Internal Error", 500)
		}
		return
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
		fmt.Sprintf("etag %q", resp.ETag)
		// see https://github.com/aws/aws-sdk-go/issues/304
		// rw.Header().Set("Etag", *resp.ETag)
	}

	if contentType != "" {
		rw.Header().Set("Content-Type", contentType)
	}
	if resp.ContentLength != nil && *resp.ContentLength > 0 {
		rw.Header().Set("Content-Length", fmt.Sprintf("%d", *resp.ContentLength))
	}
	
	io.Copy(rw, resp.Body)
	resp.Body.Close()
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
