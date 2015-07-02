s3_private_http
---------------

Private HTTP Server for Amazone S3 content.

Amazon S3 provides a public HTTP interface for accessing content, but what if you don't want publicly accessible files?

`s3_private_http` exposes a private HTTP endpoint for a s3 bucket so you can controll access to the data. This is ideal for accessing S3 via HTTP api as a backend data service, or for local http browsing of a private s3 bucket, or for use behind another authentication service (like [oauth2_proxy](https://github.com/bitly/oauth2_proxy)) to secure access.


```
Usage of ./s3_private_http:
  -bucket="": s3 bucket name
  -listen=":8080": address:port to listen on.
  -log-requests=true: log HTTP requests
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
