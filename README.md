private_s3_httpd
---------------

Private HTTP Server for Amazon S3 content.

[![Build Status](https://secure.travis-ci.org/jehiah/private_s3_httpd.png?branch=master)](http://travis-ci.org/jehiah/private_s3_httpd) [![GitHub release](https://img.shields.io/github/release/jehiah/private_s3_httpd.svg)](https://github.com/jehiah/private_s3_httpd/releases/latest)

Amazon S3 provides a public HTTP interface for accessing content, but what if you don't want publicly accessible files?

`private_s3_httpd` exposes a private HTTP endpoint for an Amazon S3 bucket so you can controll access to the data. This is ideal for accessing S3 via HTTP api as a backend data service, or for local http browsing of a private s3 bucket, or for use behind another authentication service (like [oauth2_proxy](https://github.com/bitly/oauth2_proxy)) to secure access.


```
Usage of bin/private_s3_httpd:
  -bucket string
    	S3 bucket name
  -listen string
    	address:port to listen on. (default ":8080")
  -log-requests
    	log HTTP requests (default true)
  -region string
    	AWS S3 Region (default "us-east-1")
  -s3-endpoint string
    	alternate http://address for accessing s3 (for configuring with minio.io)
```

## Configuring S3 Credentials

Before using, ensure that you've configured credentials appropriately. The best way to configure credentials is to use the `~/.aws/credentials` file, which might look like:

```
[default]
aws_access_key_id = AKID1234567890
aws_secret_access_key = MY-SECRET-KEY
```

You can learn more about the credentials file from [this blog post](http://blogs.aws.amazon.com/security/post/Tx3D6U6WSFGOK2H/A-New-and-Standardized-Way-to-Manage-Credentials-in-the-AWS-SDKs).

Alternatively, you can set the following environment variables:

```
AWS_ACCESS_KEY_ID=AKID1234567890
AWS_SECRET_ACCESS_KEY=MY-SECRET-KEY
```

