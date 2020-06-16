FROM golang:1.14.4 as builder
COPY . .
RUN go get github.com/bitly/go-simplejson
RUN go get github.com/gorilla/mux
RUN go get github.com/gorilla/handlers
RUN go build main.go

CMD ["./main"]
EXPOSE 8000
