FROM golang:1.16 as builder

WORKDIR /go/src/app
COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -tags timetzdata -o bot

FROM alpine:3.12

WORKDIR /opt

COPY --from=builder /go/src/app/bot ./bot

CMD ["/opt/bot"]
