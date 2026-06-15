# MWS Cloud Platform Terraform Provider

[![PkgGoDev](https://pkg.go.dev/badge/go.mws.cloud/terraform-provider-mws)](https://pkg.go.dev/go.mws.cloud/terraform-provider-mws)
[![Build](https://go.mws.cloud/badges/mws-cloud-platform/terraform-provider-mws)](https://go.mws.cloud/ci/mws-cloud-platform/terraform-provider-mws)
[![Go Report Card](https://goreportcard.com/badge/go.mws.cloud/terraform-provider-mws)](https://goreportcard.com/report/go.mws.cloud/terraform-provider-mws)
![Last Commit](https://img.shields.io/github/last-commit/mws-cloud-platform/terraform-provider-mws)
![Go Version](https://img.shields.io/badge/Go-1.25.6%2B-blue)

The MWS Cloud Platform Terraform Provider is a plugin that allows [Terraform](https://www.terraform.io) to manage
resources on the [MWS Cloud Platform](https://mws.ru/cloud-platform).

- [MWS Cloud Platform Terraform Provider](#mws-cloud-platform-terraform-provider)
    - [Installation](#installation)
        - [Building From Source](#building-from-source)
    - [Getting Started](#getting-started)
    - [Examples](#examples)
    - [Documentation](#documentation)
    - [Get Help](#get-help)
    - [Creators](#creators)

## Installation

To install the pre-built provider, refer to
the [tutorial](https://mws.ru/docs/cloud-platform/terraform/general/terraform-quickstart.html#provider).

### Building From Source

1. Build the provider using `go build` or `go install`.
2. Configure the `dev_overrides` block in the `~/.terraformrc` file:
    ```hcl
    provider_installation {
      dev_overrides {
        # Specify path to the provider binary here.
        "mws-cloud/mws" = "/path/to/local/provider"
      }
    }
    ```
3. Configure the provider:
    ```hcl
    terraform {
      required_providers {
        mws = {
          source = "mws-cloud/mws"
        }
      }
    }
    ```

Learn more at
the [official tutorial](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider#prepare-terraform-for-local-provider-install).

## OpenTofu Support

The MWS Cloud Platform Terraform Provider is compatible with [OpenTofu](https://opentofu.org/), an open-source alternative to Terraform. You can use this provider with OpenTofu by following the same installation and configuration steps outlined above.

## Getting Started

* [Quick Start with the MWS Cloud Platform Terraform Provider](https://mws.ru/docs/cloud-platform/terraform/general/terraform-quickstart.html)

## Examples

Check more examples in the [examples](./examples) directory.

## Documentation

* [MWS Cloud Platform Terraform Provider Documentation](https://mws.ru/docs/cloud-platform/terraform/general/whatis-terraform.html)

## Get Help

Ask for help using the [MWS Cloud Platform Support Center](https://mws.ru/docs/support/about.html).

## Creators

Created and maintained by [MWS Cloud Platform](https://mws.ru/cloud-platform).
