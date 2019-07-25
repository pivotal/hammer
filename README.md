# hammer - wrapper CLI for interacting with PCF environments
[![Build Status](https://travis-ci.com/pivotal/hammer.svg?token=jUqzM7hyJNi7CRu5xyLL&branch=master)](https://travis-ci.com/pivotal/hammer)

## Installation

The latest build of the `hammer` cli is available from the releases page.
Download the tar for your platform, untar it, and move it to your $PATH.

Alternatively you can build `hammer` from source if you have Go:
```bash
$ git clone git@github.com:pivotal/hammer.git && cd hammer && go install
```

## Config

In order to run the `hammer` tool against a given environment you need to have an environment config file in the following format:
```json
{
  "name": "ENVIRONMENT-NAME",
  "ops_manager": {
    "url": "OPSMAN-FQDN",
    "username": "OPSMAN-USERNAME",
    "password": "OPSMAN-PASSWORD"
  },
  "ops_manager_private_key": "OPSMAN-RSA-PRIVATE-KEY",
  "ops_manager_public_ip": "OPSMAN-PUBLIC-IP",
  "sys_domain": "PAS-SYSTEM-DOMAIN"
}
```
This file can then be passed into the tool via `hammer -t path-to-env-config.json <command>`.

## Development

Unit and integration tests can be run if you have [Ginkgo](https://github.com/onsi/ginkgo) installed:
```bash
$ ginkgo -r .
```

Linters can also be run using [golangci-lint](https://github.com/golangci/golangci-lint):
```bash
$ golangci-lint run
```

---

Special thanks to @blgm for letting an internal tool he created serve as the basis for this tool.
