## We specify the base image we need for our
## go application
FROM golang:1.18.1
ENV MONGOURI mongodb://mongodb:27017  
ENV REDIS_URI redis
ENV REDIS_PORT "6379"
## We create an /app directory within our
## image that will hold our application source
## files
RUN mkdir /app
## We copy everything in the root directory
## into our /app directory
ADD . /app
## We specify that we now wish to execute 
## any further commands inside our /app
## directory
WORKDIR /app
## we run go build to compile the binary
## executable of our Go program
RUN go build -o main .
## Our start command which kicks off
## our newly created binary executable

EXPOSE 4000

CMD ["/app/main"]
