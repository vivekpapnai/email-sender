FROM golang:1.16.6-alpine3.14 AS builder
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates && apk add tzdata
WORKDIR /server
ENV GO111MODULE=on
COPY go.mod /server/
COPY go.sum /server/

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=readonly -o /go/bin/create-zip

FROM scratch
COPY --from=builder /go/bin/uploadCSV /go/bin/create-zip
EXPOSE 8080
ENTRYPOINT ["/go/bin/create-zip"]