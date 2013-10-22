goprismic
=========

[![Build Status](https://secure.travis-ci.org/SoCloz/goprismic.png?branch=master)](http://travis-ci.org/SoCloz/goprismic)

Prismic.io client kit for Go

This client is still a work in progress, but should be stable - use at your own risk.

Links are currently not implemented (yet).

Usage
-----

```go
api, err := goprismic.Get("http://myrepo.prismic.io/api", "repo key")

docs, err := api.Master().Form("everything").Submit()
if err != nil {
	// handle error
}
if len(docs) == 0 {
	// nothing found
}
doc := docs[0]

st, found := doc.GetStructuredText("content")
if found {
	fmt.Println(st.AsHtml())
}
```

Documentation
-------------

See http://godoc.org/github.com/SoCloz/goprismic

License
-------

This bundle is released under the MIT license (see LICENSE).