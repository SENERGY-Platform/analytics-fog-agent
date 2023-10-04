FROM golang:1.21 AS builder

COPY . /go/src/app
WORKDIR /go/src/app

ENV GO111MODULE=on

RUN CGO_ENABLED=0 go build -o app

RUN git log -1 --oneline > version.txt

FROM docker:dind
WORKDIR /root/

COPY --from=builder /go/src/app/app .
COPY --from=builder /go/src/app/version.txt .

ENTRYPOINT ["./app"]
