# ninsho

All-inclusive compornent of sso authentication.

Note that this is only for [ochano.co proxy](https://github.com/ochanoco/proxy), so it does not support other use cases.

[![Go](https://github.com/ochanoco/ninsho/actions/workflows/go.yml/badge.svg?branch=develop)](https://github.com/ochanoco/ninsho/actions/workflows/go.yml)

## Support platform
### IdP

- **line login**
  - LINE Login is a service that allows users to log in to other apps and websites using their LINE account. 


### Web Framework

Using our supported web framework, you can use a special extension library for them. (Alternatively, you can also implement SSO using this library without supported web frameworks.)

- **gin library**
  - sso library for [gin](https://gin-gonic.com/) (Golang Web appication framework)


## How to use

Please read [examples](./example/).

To use the default configuration in the extension for gin,
you also need to set the environment variables as follows:


```sh
export NINSHO_CLIENT_ID="xxxx"
export NINSHO_CLIENT_SECRET="xxxx"

# "xxxx/callback"
export  NINSHO_REDIRECT_URI="xxxx"
```
