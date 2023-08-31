FROM alpine:3.16.0

# set up a working directory
WORKDIR /app

# copy over only the compiled binary
COPY httpServer .

# expose the port 8080 
EXPOSE 8080

# rnu the compiled binary
CMD ["./httpServer"]
