goprismic
=========

[![Build Status](https://secure.travis-ci.org/SoCloz/goprismic.png?branch=master)](http://travis-ci.org/SoCloz/goprismic)

Prismic.io client kit for Go

This client is currently used in production on http://www.socloz.com.

Go 1.3 is required.

Usage
-----

```go
// start api with the default config (3 workers, 5 req/s max, 5 seconds timeout on requests)
api, err := goprismic.Get("https://myrepo.prismic.io/api", "repo key", goprismic.DefaultConfig)

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
proxy, err := proxy.New("https://myrepo.prismic.io/api", "repo key", goprismic.DefaultConfig, proxy.Config{CacheSize: 1000, TTL: 1*time.Hour, Grace: 10*time.Minute})

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

Workers & Rate limiting
-----------------------

Access to the prismic api is done using workers, limiting the number of simultaneous connexions to the API.

A maximum number of requests per second is also set. If the threshold is passed, errors are returned until the next second.


Documentation & links
---------------------

See http://godoc.org/github.com/SoCloz/goprismic for the api documentation and http://godoc.org/github.com/SoCloz/goprismic/proxy for the proxy documentation.

Blog post : http://techblog.socloz.com/2013/11/introducing-goprismic-a-prismic-io-go-client-kit/

License
-------

This bundle is released under the MIT license (see LICENSE).