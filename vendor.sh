#!/bin/bash

gb vendor fetch -tag v0.6.4 github.com/aws/aws-sdk-go/aws
gb vendor fetch -tag v0.6.4 github.com/aws/aws-sdk-go/service/s3
gb vendor fetch -revision a98ad7ee00ec53921f08832bc06ecf7fd600e6a1 github.com/vaughan0/go-ini
gb vendor fetch --revision 8066bb491b4ea129dd71b306aab4c52868503404 github.com/gorilla/handlers
