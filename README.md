# pcf - wrapper CLI for interacting with PCF environments
[![Build Status](https://travis-ci.com/pivotal/pcf-cli.svg?token=jUqzM7hyJNi7CRu5xyLL&branch=master)](https://travis-ci.com/pivotal/pcf-cli)

## Installation

The latest build of the `pcf` cli is available from the releases page.
Download the tar for your platform, untar it, and move it to your $PATH.

Alternatively you can build `pcf` from source if you have Go:
```bash
$ git clone git@github.com:pivotal/pcf-cli.git && cd pcf-cli && go install github.com/pivotal/pcf-cli/cmd/pcf
```

## Config

In order to run the `pcf` tool against a given environment you need to create an environment config file in the following format:
```json
{
  "ert_cidr": "PAS-SUBNET-CIDR",
  "name": "ENVIRONMENT-NAME",
  "ops_manager": {
    "url": "OPSMAN-FQDN",
    "username": "OPSMAN-USERNAME",
    "password": "OPSMAN-PASSWORD"
  },
  "ops_manager_cidr": "OPSMAN-SUBNET-CIDR",
  "ops_manager_private_key": "OPSMAN-RSA-PRIVATE-KEY",
  "ops_manager_public_ip": "OPSMAN-PUBLIC-IP",
  "services_cidr": "SERVICES-SUBNET-CIDR",
  "sys_domain": "PAS-SYSTEM-DOMAIN"
}
```
This file can then be passed into the tool via `pcf -t path-to-env-config.json <command>`.

## Development

Unit and integration tests can be run if you have [Ginkgo](https://github.com/onsi/ginkgo) installed:
```bash
$ ginkgo -r .
```

Linters can also be run using [golangci-lint](https://github.com/golangci/golangci-lint):
```bash
$ golangci-lint run
```