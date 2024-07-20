# Terraform Provider Shopify

[![License: MPL-2.0](https://img.shields.io/badge/License-MPL2.0-blue.svg)](./LICENSE)
[![Tests Workflow](https://github.com/k-yomo/terraform-provider-shopify/workflows/Tests/badge.svg)](https://github.com/k-yomo/terraform-provider-shopify/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/k-yomo/terraform-provider-shopify/branch/main/graph/badge.svg)](https://codecov.io/gh/k-yomo/terraform-provider-shopify)
[![Go Report Card](https://goreportcard.com/badge/k-yomo/terraform-provider-shopify)](https://goreportcard.com/report/k-yomo/terraform-provider-shopify)

Terraform Provider for [Shopify](https://www.shopify.com/).

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Fill this in for each provider

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
