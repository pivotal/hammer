#!/usr/bin/env bash

set -eux

pushd pcf-cli
  go get github.com/onsi/ginkgo/ginkgo
  go mod download
  ginkgo -r .
popd
