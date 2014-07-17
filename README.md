goprismic
=========

[![Build Status](https://secure.travis-ci.org/SoCloz/goprismic.png?branch=master)](http://travis-ci.org/SoCloz/goprismic)

Prismic.io client kit for Go

This client should be feature complete and stable. It is currently used in production on http://www.socloz.fr.

Usage
-----

```go
// start api with 3 workers
api, err := goprismic.Get("https://myrepo.prismic.io/api", "repo key", 3)

docs, err := api.Master().Form("everything").Query("[[:d = at(document.tags, [\"Featured\"])]]").Order("my.product.name", goprismic.OrderAsc).Page(1).Submit()
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

A simple caching proxy is included.

```go
// Up to 1000 documents will be cached for up to 1 hour. Documents will be asynchronously refreshed
// if accessed 10 minutes before expiration (or later).
proxy, err := proxy.New("https://myrepo.prismic.io/api", "repo key", 3, 1000, 1*time.Hour, 10*time.Minute)

// Not cached
docs, err := proxy.Direct().Master().Form("everything").Submit()

// Cached
doc, err := proxy.GetDocument(id)

// Cached
doc, err := proxy.GetDocumentBy("product", "fieldname", "fieldvalue")

// Cached
doc, err := proxy.GetDocumentBy("product", "fieldname", "fieldvalue")

// Cached
res, err := proxy.Search().Form("menu").PageSize(200).Submit()
```

Workers
-------

Access to the prismic api is done using workers, limiting the number of simultaneous connexions to the API.


Documentation & links
---------------------

See http://godoc.org/github.com/SoCloz/goprismic for the api documentation and http://godoc.org/github.com/SoCloz/goprismic/proxy for the proxy documentation.

Blog post : http://techblog.socloz.com/2013/11/introducing-goprismic-a-prismic-io-go-client-kit/

License
-------

This bundle is released under the MIT license (see LICENSE).