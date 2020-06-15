FROM golang:1.10.1-alpine3.7 as builder
COPY main.go .
RUN go get github.com/bitly/go-simplejson
RUN go get github.com/gorilla/mux
RUN go get github.com/gorilla/handlers
RUN go build -o /main main.go

FROM alpine:3.7  
CMD ["./main"]
COPY --from=builder /main .
