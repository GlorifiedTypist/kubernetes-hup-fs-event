FROM golang:alpine
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o hup-fs-event -v .

FROM alpine
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true
COPY --from=0 /go/src/app/hup-fs-event .
ENV WATCH_DIRECTORY="/"
ENV HUP_TARGET="hup-fs-event"
ENTRYPOINT ["/hup-fs-event"]