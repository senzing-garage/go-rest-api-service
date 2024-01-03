# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
[markdownlint](https://dlaa.me/markdownlint/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.1] - 2023-11-02

### Changed in 0.1.1

- Update dependencies
  - github.com/ogen-go/ogen v0.77.0
  - github.com/senzing-garage/go-common v0.3.2-0.20231018174900-c1895fb44c30
  - github.com/senzing-garage/go-sdk-abstract-factory v0.4.3

## [0.1.0] - 2023-10-25

### Changed in 0.1.0

- Supports SenzingAPI 3.8.0
- Refactor to [template-go](https://github.com/senzing-garage/template-go)
- Deprecated functions have been removed
- Update dependencies
  - github.com/ogen-go/ogen v0.76.0
  - github.com/senzing-garage/g2-sdk-go v0.7.4
  - github.com/senzing-garage/go-common v0.3.1
  - github.com/senzing-garage/go-logging v1.3.3
  - github.com/senzing-garage/go-observing v0.2.8
  - github.com/senzing-garage/go-sdk-abstract-factory v0.4.2
  - go.opentelemetry.io/otel v1.19.0
  - go.opentelemetry.io/otel/metric v1.19.0
  - go.opentelemetry.io/otel/trace v1.19.0
  - google.golang.org/grpc v1.59.0

## [0.0.6] - 2023-09-01

### Changed in 0.0.6

- Last version before SenzingAPI 3.8.0

## [0.0.5] - 2023-08-07

### Changed in 0.0.5

- Refactor to `template-go`
- Update dependencies
  - github.com/go-faster/jx v1.1.0
  - github.com/ogen-go/ogen v0.72.1
  - github.com/senzing-garage/g2-sdk-go v0.6.8
  - github.com/senzing-garage/go-common v0.2.11
  - github.com/senzing-garage/go-logging v1.3.2
  - github.com/senzing-garage/go-observing v0.2.7
  - google.golang.org/grpc v1.57.0

## [0.0.4] - 2023-07-26

### Changed in 0.0.4

- Re-generate files in `senzingrestapi` package
- Update dependencies
  - github.com/go-faster/jx v1.0.1
  - github.com/ogen-go/ogen v0.72.0
  - github.com/senzing-garage/go-common v0.2.5

## [0.0.3] - 2023-07-14

### Changed in 0.0.3

- Update dependencies
  - github.com/ogen-go/ogen v0.71.0
  - github.com/senzing-garage/g2-sdk-go v0.6.7
  - github.com/senzing-garage/g2-sdk-json-type-definition v0.1.1
  - github.com/senzing-garage/go-common v0.2.1
  - github.com/senzing-garage/go-logging v1.3.1
  - github.com/senzing-garage/go-sdk-abstract-factory v0.3.1
  - google.golang.org/grpc v1.56.2

## [0.0.2] - 2023-06-16

### Added to 0.0.2

- Update dependencies
  - github.com/ogen-go/ogen v0.69.1
  - github.com/senzing-garage/g2-sdk-go v0.6.6
  - github.com/senzing-garage/go-common v0.1.4
  - github.com/senzing-garage/go-logging v1.2.6
  - github.com/senzing-garage/go-observing v0.2.6
  - google.golang.org/grpc v1.56.0

## [0.0.1] - 2023-06-09

### Added to 0.0.1

- Methods for GET: `/heartbeat`, `/license`, `/version`, and `/specifications/open-api`
