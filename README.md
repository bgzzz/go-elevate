# go-elevate
microservice that gets ground height based on lat and lon coordinates

## API

coordinates are in radians 

1. /height?lon=86.922623&lat=27.986065 for single coord request 
```
$ curl '127.0.0.1:1323/height?lat=27.986065&lon=86.22623'
```

2. /heights
```
curl -X POST -H "Content-Type: application/json" \
-d  '{"coords":[{"lon":138.72905,"lat":35.360638},{"lat":27.986065,"lon":86.922623}]}' \
'127.0.0.1:1323/heights'
```

## How to compile 

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

1. unit and integration tests
1. config refactor 
1. HTTP doc for API 
1. caching of the request (redis) or may be even save in heights in regular DB (AWS tiles might be shut down)
1. support of different coordinate types (radians, degrees ... etc.) 
1. differrent algos for height calculation (median, avg, supported (bases on radians to EPSG:3857 convertion error) )
