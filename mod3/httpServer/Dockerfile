FROM golang:1.15-alpine AS build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
        GOPROXY="https://goproxy.cn,direct"

WORKDIR /go/src/httpServer

ADD main.go .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o httpServer main.go

FROM scratch AS prod

COPY --from=build /go/src/httpServer/httpServer .

ENV GIN_MODE=release \
    PORT=8080

EXPOSE 8080

CMD ["./httpServer"]
~
