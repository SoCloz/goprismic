goprismic
=========

[![Build Status](https://secure.travis-ci.org/SoCloz/goprismic.png?branch=master)](http://travis-ci.org/SoCloz/goprismic)

Prismic.io client kit for Go

This client should be feature complete and stable - however use it at your own risk.

Usage
-----

```go
api, err := goprismic.Get("https://myrepo.prismic.io/api", "repo key")

docs, err := api.Master().Form("everything").Query("[[:d = at(document.tags, [\"Featured\"])]]").Submit()
if err != nil {
	// handle error
}
if len(docs) == 0 {
	// nothing found
}
doc := docs[0]

st, found := doc.GetStructuredTextFragment("content")
if found {
	fmt.Println(st.AsHtml())
}
```

Links
-----

You have to resolve document links using a user-supplied link resolver :

```go
r := func(l link.Link) string {
	return l.(*link.DocumentLink).Document.Slug
}
```
and resolve links at document/fragment/block level :
```go
doc.ResolveLinks(r)

st, _ := doc.GetStructuredTextFragment("content")
st.ResolveLinks(r)

p, _ := st.GetFirstParagraph()
p.ResolveLinks(r)
```

Proxy
-----

A simple caching proxy is included. Only single document accesses are cached.

```go
proxy, err := proxy.New("https://myrepo.prismic.io/api", "repo key", 1*time.Hour, 10*time.Minute)

docs, err := proxy.Direct().Master().Form("everything").Submit()

...

doc, err := proxy.GetDocument(id)

...

doc, err := proxy.GetDocumentBy("product", "fieldname", "fieldvalue")
```

Documentation
-------------

See http://godoc.org/github.com/SoCloz/goprismic

License
-------

This bundle is released under the MIT license (see LICENSE).