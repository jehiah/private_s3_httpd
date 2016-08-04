#!/bin/bash

if [ -e vendor ]; then
   echo "vendor directory exists. remove before running"
   exit 1
fi

gb vendor fetch -no-recurse -tag v1.3.0 github.com/aws/aws-sdk-go
gb vendor fetch -no-recurse -revision cf53f9204df4fbdd7ec4164b57fa6184ba168292 github.com/go-ini/ini
gb vendor fetch -no-recurse -revision bd40a432e4c76585ef6b72d3fd96fb9b6dc7b68d github.com/jmespath/go-jmespath
gb vendor fetch -no-recurse -revision a98ad7ee00ec53921f08832bc06ecf7fd600e6a1 github.com/vaughan0/go-ini
gb vendor fetch -no-recurse -revision 8066bb491b4ea129dd71b306aab4c52868503404 github.com/gorilla/handlers
