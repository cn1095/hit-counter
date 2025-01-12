FROM golang:1.16.3-alpine AS builder

WORKDIR /hit-counter

RUN go env -w GO111MODULE="on"

# copy go.mod go.sum
COPY ./go.mod ./go.sum ./

# download Library
RUN go mod download

# copy all
COPY ./ ./

RUN CGO_ENABLED=0 go build -a -ldflags "-w -s" -o /hit-counter/hit-counter

# Minimize a docker image
FROM gcr.io/distroless/base:latest

COPY --from=builder /hit-counter/hit-counter /hit-counter/hit-counter

COPY /public /public

COPY /view /hit-counter/view

CMD ["/hit-counter/hit-counter"]
