

FROM golang:1.18.1

ENV TZ=Asia/Bangkok
ENV MONGOURI mongodb://mongodb:27017  
ENV REDIS_URI redis
ENV REDIS_PORT "6379"

RUN mkdir /app
COPY . /app

WORKDIR /app
RUN go build -o main .

EXPOSE 4000
CMD ["/app/main"]
