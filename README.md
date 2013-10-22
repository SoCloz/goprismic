goprismic
=========

[![Build Status](https://secure.travis-ci.org/SoCloz/goprismic.png?branch=master)](http://travis-ci.org/SoCloz/goprismic)

Prismic.io client kit for Go

This client has not (yet) been field tested - use at your own risk.

TODO :

* links

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

License
-------

This bundle is released under the MIT license (see LICENSE).