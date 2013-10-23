goprismic
=========

[![Build Status](https://secure.travis-ci.org/SoCloz/goprismic.png?branch=master)](http://travis-ci.org/SoCloz/goprismic)

Prismic.io client kit for Go

This client should be feature complete and stable - use at your own risk.

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

st, found := doc.GetStructuredText("content")
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

st, _ := doc.GetStructuredText("content")
st.ResolveLinks(r)

p, _ := doc.GetFirstParagraph()
p.ResolveLinks(r)
```

Documentation
-------------

See http://godoc.org/github.com/SoCloz/goprismic

License
-------

This bundle is released under the MIT license (see LICENSE).