# Kitty key value store

A key value store with http RESTful API.

## Running the server
``` bash
$ ./example 
kitty: 2019/02/24 15:54:06 Kitty key value server is starting ( try: http://localhost:8080/val ) ...
...
```

## Query the server

``` bash
$ # insert a new key value pair.
$ curl -s http://localhost:8080/val --data "{\"hello\": \"world\"}" | jq
{
  "hello": "world"
}
$ # insert a new key value pair.
$ curl -s http://localhost:8080/val --data "{\"kitty\": \"cat\", \"gorilla\": 123}" | jq
{
  "gorilla": 123,
  "kitty": "cat"
}
$ # query all key value pairs.
$ curl -s http://localhost:8080/val | jq
{
  "gorilla": 123,
  "hello": "world",
  "kitty": "cat"
}
$ # query a key value pair with key = ``kitty`.
$ curl -s http://localhost:8080/val/kitty | jq
{
  "kitty": "cat"
}
$ # delete a key value pair with key = ``hello`.
$ curl -s -X DELETE http://localhost:8080/val/hello | jq
{
  "hello": "world"
}
$ # query all key value pairs.
$ curl -s http://localhost:8080/val | jq
{
  "gorilla": 123,
  "kitty": "cat"
}
$ # delete a key value pair with key = `dog`, no such pair in the data store.
$ curl -s -X DELETE http://localhost:8080/val/dog | jq
{
  "error": "can't find key dog"
}

```
