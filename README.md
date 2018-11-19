# Image Server

This server can read pictures in png (first priority) and jpg (second priority) format from disk and serve them over http
as png or jpg images at the resolution requested by the client.
It can also serve raw (unmodified) versions of png, jpg, svg and pdf files.

## URL Schema and source image logic 
* /<height>p/path/to/image.{jpg,png}
  * returns 200 and a scaled version of the image to <height> pixel height if {$IMAGE_DIR}path/to/image.png is readable
  * returns 200 and a scaled version of the image to <height> pixel height if {$IMAGE_DIR}path/to/image.jpg is readable
  * returns 404 otherwise 

* /<width>w/path/to/image.{jpg,png}
  * returns 200 and a scaled version of the image to <width> pixel width if {$IMAGE_DIR}path/to/image.png is readable
  * returns 200 and a scaled version of the image to <width> pixel width if {$IMAGE_DIR}path/to/image.jpg is readable
  * returns 404 otherwise 
  
* /raw/path/to/image.<ext> (where <ext> can be one of png, jpg, svg, pdf)
  * returns 200 and the unmodified file if {$IMAGE_DIR}path/to/image.<ext> is readable
  * returns 404 otherwise

## Configuration
Configuration is done exclusively using environment variables. The following variables are supported:
* **PORT**
  * port to listen on
  * defaults to 80
  * gin overwrites this env variable automatically
* **BIND**
  * addresses to bind to
  * defaults to 127.0.0.1
* **IMAGE_DIR**
  * where to look for the input images
  * use an absolute path or a path relative to `pwd` when you execute the server.
  * Use a / at the end
* **MAX_AGE**
  * Cache-Control-Header: controls how many seconds the output is cached by reverse http proxies
  * (defaults to 0 -> no cache)
* **JPG_QUALITY**
  * quality of output jpg images 1...100
  * defaults to 90
  * this setting has no effect on /raw output

## Caching
The server does not do any internal caching and reads and scales the image on each request for which
a body has to be sent.

This server sends `Last-Modified` and `Cache-Control: "public, max-age=MAX_AGE` headers in responses.
`304 Not Modified` responses are sent when a correct `If-Modified-Since` request header is present.

As `Last-Modified` time the time of the the first matching source image is used. 
  
This results in the Browser Cache being used when directly accessing the server.

When configuring nginx as reverse proxy as shown in docker-compose.yaml and nginx/default.conf,
this results in a very high performance system.
nginx uses up to 4GB of disk space as a cache and only checks every 15s with the imgsrv if the file was modified.
The actual image is only scaled once after startup unless the nginx cache overflows.

## Deployment options

### standalone no cache
The serve can be used standalone with the following limitations:
* there is no server-side cache
* no https

### using nginx as caching reverse proxy
For production, a reverse proxy should be added for caching and possibly for ssl or http2 termination.
See docker-compose.yaml for an example setup using nginx.

## Development
For development on this server, you might want to use gin, which automatically recompiles the server
on each request if the source has changed.

```
go get github.com/codegangsta/gin
export IMAGE_DIR=/home/lk/tmp/pics_raw/
gin
```

To build the docker images use:
```
docker-compose build
```
