# Kitty key value store server

In memory key value strore http server.

## Running the server
``` bash
$ ./example 
kitty: 2019/02/24 15:54:06 Kitty key value server is starting ( try: http://localhost:8080/val ) ...
...
```

## Query the server

```
$ curl -s http://localhost:8080/val --data "{\"hello\": \"world\"}" | jq
{
  "hello": "world"
}
$ curl -s http://localhost:8080/val --data "{\"kitty\": \"cat\"}" | jq
{
  "kitty": "cat"
}
$ curl -s http://localhost:8080/val | jq
{
  "hello": "world",
  "kitty": "cat"
}
$ curl -s http://localhost:8080/val/kitty | jq
{
  "kitty": "cat"
}
$ curl -s -X DELETE http://localhost:8080/val/hello | jq
{
  "hello": "world"
}
$ curl -s http://localhost:8080/val | jq
{
  "kitty": "cat"
}
$ curl -s -X DELETE http://localhost:8080/val/cat | jq
{
  "error": "can't find key cat"
}

```