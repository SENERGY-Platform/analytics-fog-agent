FROM golang
COPY . /go/src/fog_agent
WORKDIR /go/src/fog_agent
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o fog_agent .

FROM docker:dind
WORKDIR /root/
COPY --from=0 /go/src/fog_agent/fog_agent .
RUN ls -la
CMD ./fog_agent
