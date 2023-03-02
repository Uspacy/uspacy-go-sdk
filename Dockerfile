FROM golang:1.19-alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh make

ENV SERVICE_NAME app
ENV APP /src/${SERVICE_NAME}/
ENV WORKDIR ${GOPATH}${APP}

WORKDIR $WORKDIR

ADD . $WORKDIR

RUN go get ./...
RUN go get -u golang.org/x/lint/golint

RUN go mod tidy
RUN CGO_ENABLED=0 go build -i -v -o release/app

FROM alpine

COPY --from=builder /go/src/app/release/app /

RUN chmod +x /app

ENTRYPOINT ["/app"]