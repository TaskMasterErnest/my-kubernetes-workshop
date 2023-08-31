FROM alpine:3.16.0

# create a working directory
WORKDIR /app

# copy over the main.go file
COPY memconsume .

# expose the port 8085
EXPOSE 8085 

# run the application
CMD ["./memconsume"]
