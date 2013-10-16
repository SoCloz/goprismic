goprismic
=========

Prismic.io client kit for Go

For now, this client is missing many features.

Usage
-----

```go
api, err := goprismic.Get("http://myrepo.prismic.io/api", "repo key")

docs, err := api.Master().Form("everything").Submit()
```

