![](https://img.shields.io/badge/language-go-blue.svg) ![](https://github.com/donng/go-shortenurl/workflows/go%20build/badge.svg)

`go-shortenurl` is a url shorten service，based on a free course of mooc - [Go开发短地址服务](https://www.imooc.com/learn/1150) .

## Installation

go-shortenurl requires the Go version with [Module](https://github.com/golang/go/wiki/Modules) support.

```
git clone https://github.com/donng/go-shortenurl.git

cd go-shortenurl

go mod download
```

## API

there are three simple apis

- shorten url
- get shorten url info
- visit short url and redirect

### shorten url

```
API：/api/shorten
METHOD：POST
PARAMS: { "url": "https:www.example.com", "expire_in_minutes": 60 }
```

### get shorten url info

```
API: /api/info/{link}
METHOD: GET
```

### visit short url and redirect

visit link will return status code 307 and redirect to the origin url

```
API: /{link}
METHOD: GET
```
