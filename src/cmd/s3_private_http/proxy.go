package main

import (
	"net/http"
	"io"


	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Proxy struct {
	Bucket string
	Svc *s3.S3
}

func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	resp, err := p.Svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(p.Bucket),
		Key:    aws.String(req.URL.Path),
	})
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
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
