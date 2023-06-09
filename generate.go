package main

// For ogen options, see https://github.com/ogen-go/ogen/blob/main/cmd/ogen/main.go
//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --target senzingrestapi --package senzingrestapi  --generate-tests --clean --debug.ignoreNotImplemented "sum type parameter, discriminator inference" restapiservice/openapi.yaml
