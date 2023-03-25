FROM golang:1.19.1-alpine as build
# RUN apk add make
RUN mkdir -p /app/src
WORKDIR /app/src

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o main ./cmd/main.go
# RUN make build

FROM alpine:3.17
COPY --from=mwader/static-ffmpeg:5.1.2 /ffmpeg /usr/local/bin/
COPY --from=mwader/static-ffmpeg:5.1.2 /ffprobe /usr/local/bin/

COPY --from=build /app/src/main /root/main
COPY config.yaml /root/config.yaml

WORKDIR /root

CMD ["./main"]