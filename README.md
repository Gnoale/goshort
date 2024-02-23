# URL Shortener

## Build / Run
`docker compose up`

`go run ./cmd/main.go`

## Design

Complete URLS are stored in the database, the DB schema is minimalist.

1- User submit a new URL

    - The service try to insert the URL in DB.
    - If the URL already exist, return HTTP 409 (conflict)
    - if the URL is successfully inserted, we return the SLUG

SLUG is computed each time an URL is inserted or retrieved.

It is the URL ID converted to a base62 representation string.

`0123456789ABCDEFGHIJKLMNOPQRSTUVWXYSabcdefghijklmnopqrstuvwxyz`

We assume a maximum length of 7 character for the slug part, which leave rooms for up to 62^7 different URL in database... 

2- User fetch a slug from the service

    - the service convert the slug in the base10 integer representation
    - the URL is retrieved by its ID
    - if it is found, the service return an HTTP 301 response with the original location in corresponding the header

### API

- `POST /api/v1/shorten` is the endpoint to submit a valid URL to be shorten

`curl -v --data '{"long_url":"https://github.com/mattn/go-sqlite3/blob/master/_example/simple/"}' http://127.0.0.1:8000/api/v1/shorten`

```
< HTTP/1.1 201 Created
< Date: Fri, 23 Feb 2024 15:52:25 GMT
< Content-Length: 31
< Content-Type: text/plain; charset=utf-8
< 
http://127.0.0.1:8000/api/v1/1
``` 

Returns a slug suitable to retrieve later the same URL from the service

`curl -v http://127.0.0.1:8000/api/v1/1`

```
< HTTP/1.1 301 Moved Permanently
< Location: https://github.com/mattn/go-sqlite3/blob/master/_example/simple/
```


