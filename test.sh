#!/bin/bash
set -e

gb test -timeout 60s
# gb test -timeout 60s -race
gb build
