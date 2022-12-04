# line-login

All-inclusive compornent of LINE Login authentication.

This project include as follows:

core library
: core sso library for line login

gin library
: sso library for [gin](https://gin-gonic.com/) (Golang Web appication framework)


## How to use

Please read [examples](./example/).

## Note

Some parameters are needed to set as environment variables.

```sh
export CLIENT_ID="xxx"
export CLIENT_SECRET="xxx"

# "xxx/callback"
export REDIRECT_URI="xxx"
```