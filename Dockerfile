FROM golang:alpine AS builder
ENV CGO_ENABLED=0
RUN apk update && apk add --no-cache git
#WORKDIR $GOPATH/src/mypackage/myapp/
WORKDIR /app
COPY . .
RUN go mod tidy
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/app /app/cmd/web


FROM scratch
EXPOSE 4000
COPY --from=builder /go/bin/app /go/bin/app
ENTRYPOINT ["/go/bin/application"]