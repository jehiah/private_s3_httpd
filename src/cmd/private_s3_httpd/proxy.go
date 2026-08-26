package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

type Proxy struct {
	Bucket string
	Svc    *s3.Client
}

func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	key := strings.TrimLeft(req.URL.Path, "/")
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

	var is304 bool
	resp, err := p.Svc.GetObject(req.Context(), input)
	if err != nil {
		var noSuchKey *types.NoSuchKey
		var respErr *smithyhttp.ResponseError
		switch {
		case errors.As(err, &noSuchKey):
			http.Error(rw, "Page Not Found", 404)
			return
		case errors.As(err, &respErr) && respErr.HTTPStatusCode() == http.StatusNotModified:
			is304 = true
			// continue so other headers get set appropriately
		default:
			log.Printf("Error: %v", err)
			http.Error(rw, "Internal Error", 500)
			return
		}
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
		if _, err := io.Copy(rw, resp.Body); err != nil {
			log.Printf("Error copying response body: %v", err)
		}
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}
}
