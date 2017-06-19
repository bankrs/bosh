# bosh - an interactive shell for bankrs OS 

This is an interactive shell for accessing the Bankrs OS API.

**Documentation:** [![GoDoc](https://godoc.org/github.com/bankrs/bosh?status.svg)](https://godoc.org/github.com/bankrs/bosh)  

bosh requires Go version 1.7 or greater.

## Getting started

Ensure you have a working Go installation and then use go get as follows:

```
go get github.com/bankrs/bosh
```

## Running bosh

Running `bosh` should start the interactive shell, assuming $GOPATH/bin is in your path.

Type `help` to get a list of commands.

## Example: searching financial providers

Login with a developer account and use the assigned application ID:

```
> login email@example.com
Password: *******
> useapp df4ef6c1-f12c-40ec-826e-c049874763de
df4ef6c1-f12c-40ec-826e-c049874763de> searchproviders deutsch
[
  {
    "score": 1,
    "provider": {
      "id": "DE-BIN-12030000",
      "name": "Deutsche Kreditbank Berlin",
      "description": "",
      "country": "DE",
      "url": "",
      "address": "10117 Berlin",
      "postal_code": "10117",
      "challenges": [
        {
          "id": "login",
          "desc": "Legitimations-ID/Anmeldename",
          "type": "alphanumeric",
          "secure": false,
          "unstoreable": false
        },
        {
          "id": "pin",
          "desc": "Onlinebanking-PIN",
          "type": "alphanumeric",
          "secure": true,
          "unstoreable": true
        }
      ]
    }
  }
]
```
