FROM golang:1.16 as builder
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 make build

FROM gcr.io/distroless/static-debian10:latest-amd64
WORKDIR /opt
COPY --from=builder /go/src/app/tkkz-bot ./bot
ENTRYPOINT ["/opt/bot"]
