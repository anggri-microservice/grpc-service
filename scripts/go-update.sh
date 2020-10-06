#!/bin/bash
export GOPRIVATE=github.com/anggri-microservice/*

cd ./cmd/grpc-service || exit 1
go get -u
go mod vendor
