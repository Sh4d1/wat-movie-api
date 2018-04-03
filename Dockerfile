FROM golang:1.10 as builder

WORKDIR /go/src/github.com/Sh4d1/wat-movie-api

COPY . .

#RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app

COPY --from=builder /go/src/github.com/Sh4d1/wat-movie-api .

CMD ["./wat-movie-api"]
