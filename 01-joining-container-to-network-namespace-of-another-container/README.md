# Joining Two Container Network via the Network Namespace

## CheckList
 [] Create a simple HTTP server in Go
 [] Create a Dockerfile for the Go application
 [] Tag and push the image built to an image repository
 [] Joining another container to running HTTP server container
 [] Conclusion


## Create a simple HTTP server in Go
- First we create a simple application in Go that brings up the words, "Welcome to my Kubernetes Workshop!", when ran.

```Go
package main

import (
  "fmt"
  "log"
  "net/http"
)

func main() {
  http.HandleFunc("/", handler)
  log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
  log.Printf("Ping from %s", r.RemoteAddr)
  fmt.Println(w, "Welcome to my Kubernetes Workshop!")
}
```

## Create a Dockerfile for the Go application
- With this Dockerfile, Go is a compiled language so we can use Alpine as a base image.
- In this, we compile the simple HTTP server application first with the command `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o httpServer main.go`. It can now run anywhere, no dependencies needed.
- We construct a simple Dockerfile to take advantage of the compiled binary.

```Dockerfile
FROM alpine:3.16.8

# set up a working directory
WORKDIR /app 

# copy over the compiled binary only
COPY httpServer . 

# expose the port 8080
EXPOSE 8080 

# run the application
CMD ["./httpServer"]
```

