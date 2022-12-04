# line-login

All-inclusive compornent of LINE Login authentication.

This project include as follows:

- **core library**
  - core sso library for line login

- **gin library**
  - sso library for [gin](https://gin-gonic.com/) (Golang Web appication framework)


## How to use

Please read [examples](./example/).

## Note

Some parameters are needed to set as environment variables.

```sh
provider.ClientID = os.Getenv("CLIENT_ID")
provider.ClientSecret = os.Getenv("CLIENT_SECRET")
provider.RedirectURL = os.Getenv("REDIRECT_URL") + callback
```
