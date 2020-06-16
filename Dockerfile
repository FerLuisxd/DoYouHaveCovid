FROM golang:1.14.4 as builder
COPY go.mod go.sum ./
RUN go mod download
RUN go get github.com/bitly/go-simplejson
RUN go get github.com/gorilla/mux
RUN go get github.com/gorilla/handlers
COPY main.go .
RUN go build -o /main main.go

FROM alpine:3.7  
EXPOSE 8000
CMD ["./main"]
COPY --from=builder /main .
