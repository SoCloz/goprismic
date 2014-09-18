goprismic
=========

[![Build Status](https://secure.travis-ci.org/SoCloz/goprismic.png?branch=master)](http://travis-ci.org/SoCloz/goprismic)

Prismic.io client kit for Go

This client is currently used in production on http://www.socloz.com.

Go 1.3 is required.

Usage
-----

```go
// start api with the default config (3 workers, 5 seconds timeout on requests)
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

A simple asynchronous caching proxy is included.

```go
// Up to 1000 documents will be cached.
proxy, err := proxy.New("https://myrepo.prismic.io/api", "repo key", goprismic.DefaultConfig, proxy.Config{CacheSize: 1000})

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

When a document is updated :
* between "update time" and "update time + TTL", documents are asynchronously refreshed (if a document is accessed, the cached version is returned, and the cache is asynchronously updated),
* after "update time + TTL", documents are fetched directly from prismic (cache miss).

If no ttl is set, refreshes are always asynchronous.

You can define :
* a refresh chance (between 0 and 1) - only a fraction of contents will be refreshed at a time, ensuring that prismic is not flooded after an update,
* a master refresh interval - the proxy will check for updates at the defined frequency.

The proxy will try to avoid flooding prismic by automatically lower/raise the refresh chance.

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