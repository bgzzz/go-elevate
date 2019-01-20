# go-elevate
microservice that gets ground height based on lat and lon coordinates

## API

coordinates are in radians 

1. GET /heights?coords=%5B%7B%22lon%22%3A138.72905%2C%22lat%22%3A35.360638%7D%2C%7B%22lat%22%3A27.986065%2C%22lon%22%3A86.922623%7D%5D
```
$ curl -G -v "http://localhost:1323/heights" --data-urlencode 'coords=[{"lon":138.72905,"lat":35.360638},{"lat":27.986065,"lon":86.922623}]'
```

2. POST /heights
```
curl -X POST -H "Content-Type: application/json" \
-d  '{"coords":[{"lon":138.72905,"lat":35.360638},{"lat":27.986065,"lon":86.922623}]}' \
'127.0.0.1:1323/heights'
```

## How to compile 

Resolver dependencies first

```
$ glide install
```

And then compile the code

``` 
$ make
```

## How to build docker 

```
$ make dock-build
```

## How to docker compose 

```
$ make dock-compose
``` 

## TODO

1. HTTP doc for API 
1. caching of the request (redis) or may be even save in heights in regular DB (AWS tiles might be shut down)
1. support of different coordinate types (radians, degrees ... etc.) 
