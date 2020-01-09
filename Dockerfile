FROM golang:alpine AS builder
ENV CGO_ENABLED=0
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .
RUN go get -d -v
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/application


FROM scratch
EXPOSE 4000
COPY --from=builder /go/bin/application /go/bin/application
ENTRYPOINT ["/go/bin/application"]