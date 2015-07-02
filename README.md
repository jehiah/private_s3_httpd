s3_private_http
---------------

Amazon S3 provides a HTTP interface for accessing content, but what if you don't want publicly accessible files?

`s3_private_http` exposes a private s3 bucket via HTTP so you can controll access to the data? This is idea for exposing a HTTP api as a backend data service, or for local browsing of a s3 bucket, or for using behind another service (like [oauth2_proxy](https://github.com/bitly/oauth2_proxy)) to secure content.


