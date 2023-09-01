# Connecting Two Containers Together

## CheckList
- [] Write an application that counts the number of times a page is visited
- [] Rewrite application to connect to a Redis datastore backend
- [] Test the connectivity

## Write a simple Page View application
- Here, we write a simple application that prints out "Welcome! You are visitor #" anytime the page is visited.
- The code for this is:

```Go
package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    log.Printf("Ping from %s", r.RemoteAddr)
    // Increment the number of visitors.
    visitorCount++

    // Print the number of visitors to the ResponseWriter.
    fmt.Fprintf(w, "Welcome! You are visitor #%d\n", visitorCount)
}

func main() {
    // Create a new http.Server and listen on port 8080.
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
```

## Connecting to a Redis database backend
- Here, we want to persist the visitorCount number in the datastore so that even if the container is restarted, it will continue working as it was.
- We have to modify the code in order to handle the datastore connection we are introducing, the final code is in the `main.go` file.
- For Podman, do the following things:
  - create a Pod and publish the ports 8080 and 6379 on these pods; `podman pod create --name pgView-srv --publish 8080 --publish 6397`.
  - start the Redis container in this new Pod with `podman run -d --name db --pod pgView-srv redis:alpine3.18`.
  - start the page view application container also with `podman run -d --name pgview --pod pgView-srv ernestklu/pageviewer:v0.0.1`.
  - since the container is minimal and we will like to query the page view application to find the visitor counts, we have to use the nicolaka/netshoot application again. Inject it into the page viewer application with the command `podman run -it --pod pgView-srv --net container:29359ba4dd06 nicolaka/netshoot:v0.11`. The ID should be replaced by the ID of the pgview container.


## Test the connectivity
- With the tool now ready, we can get into the pgview container and call the application with `curl localhost:8080`.
- It should return a message saying 'Welcome! You are visitor #1'. Do this as many times as you like and note the visitor count number that was the last.
- Then restart the pgview container with `podman restart pgview`.
- Run the test with curl again and it should have incremented on the last value you saw before restarting the container.

## WHY?
- This is to illustrate the concept of container orchestration and the problems it aims to solve.
