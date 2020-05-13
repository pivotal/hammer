# hammer - wrapper CLI for interacting with PCF environments
[![Test Status](https://github.com/pivotal/hammer/workflows/Test/badge.svg)](https://github.com/pivotal/hammer/actions)

## Installation

The latest build of the `hammer` cli is available from the releases page.
Download the tar for your platform, untar it, and move it to your $PATH.

Or using `brew` on macOS or Linux:
```bash
brew tap pivotal/hammer https://github.com/pivotal/hammer
brew install hammer
```

Alternatively you can build `hammer` from source if you have Go installed:
```bash
git clone git@github.com:pivotal/hammer.git && cd hammer && go install
```

## Config

In order to run the `hammer` tool against a given environment you need to have an environment config file in the following format:
```json
{
  "name": "ENVIRONMENT-NAME",
  "ops_manager": {
    "url": "OPSMAN-URL",
    "client_id": "OPSMAN-CLIENT-ID",
    "client_secret": "OPSMAN-CLIENT-SECRET",
    "username": "OPSMAN-USERNAME",
    "password": "OPSMAN-PASSWORD"
  },
  "ops_manager_private_key": "OPSMAN-RSA-PRIVATE-KEY",
  "ops_manager_public_ip": "OPSMAN-PUBLIC-IP",
  "ops_manager_ssh_user": "OPSMAN-SSH-USER",
  "sys_domain": "PAS-SYSTEM-DOMAIN",
  "pks_api":  {
     "url": "PKS-API-URL"
  }
}
```
Or the equivalent in yaml:
```yaml
name: ENVIRONMENT-NAME
ops_manager:
  client_id: OPSMAN-CLIENT-ID
  client_secret: OPSMAN-CLIENT-SECRET
  password: OPSMAN-PASSWORD
  url: OPSMAN-URL
  username: OPSMAN-USERNAME
ops_manager_private_key: OPSMAN-RSA-PRIVATE-KEY
ops_manager_public_ip: OPSMAN-PUBLIC-IP
ops_manager_ssh_user: OPSMAN-SSH-USER
pks_api:
  url: PKS-API-URL
sys_domain: PAS-SYSTEM-DOMAIN
```
This file can then be passed into the tool via `hammer -t path-to-env-config <command>`.

`ops_manager_ssh_user` is an optional field and if not set then `hammer -t path-to-env-config ssh opsman` will use
`ubuntu` to ssh to the OpsManager VM, if users need to ssh via a different username they should set this as appropriate.

Only one set of `ops_manager.client_id` and `ops_manager.client_secret` or `ops_manager.username` and `ops_manager.password`
need to be specified in order to use `hammer om`, if both sets are specified then in line with `om` the client details
will be used.

NB: `sys_domain` and `pks_api.url` are only needed for using `hammer cf-login` and `hammer pks-login` respectively.

## Development

Unit and integration tests can be run if you have [Ginkgo](https://github.com/onsi/ginkgo) installed:
```bash
ginkgo -r .
```

Linters can also be run using [golangci-lint](https://github.com/golangci/golangci-lint):
```bash
golangci-lint run
```

Or just run both with:
```
make test
```

---

Special thanks to [@blgm](https://github.com/blgm) for letting an internal tool he created serve as the basis for this tool.
