FROM golang:alpine as builder

ARG CGO_ENABLED=1
ARG CC=gcc

RUN mkdir /build /out && apk add --no-cache git ca-certificates gcc g++

WORKDIR /build
COPY . .
RUN go get && go build -o /out/dynqr


FROM alpine:latest

RUN mkdir /app /data && addgroup -S dynqr && adduser -S dynqr -G dynqr && apk add --no-cache ca-certificates
COPY --from=builder --chown=dynqr:dynqr /out/dynqr /app/dynqr
COPY --chown=dynqr:dynqr static/ /app/static
RUN chown -R dynqr:dynqr /app/ /data/ && chmod -R 774 /app/ /data/
USER dynqr

ENTRYPOINT /app/dynqr

