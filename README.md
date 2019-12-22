# ShortLinks

This is a little URL shortener service that is primarily intended for
personal use by a single person. It stores all its data in a local
SQLite database.

## Usage

```
$ shortlinks --addr localhost:8000 &

$ curl -X POST http://localhost:8000 --form url=http://zerokspot.com
17l3x

$ curl -i http://localhost:8000/17l3x
HTTP/1.1 307 Temporary Redirect
Content-Type: text/html; charset=utf-8
Location: http://zerokspot.com
Date: Sun, 22 Dec 2019 10:29:56 GMT
Content-Length: 56

<a href="http://zerokspot.com">Temporary Redirect</a>.
```
