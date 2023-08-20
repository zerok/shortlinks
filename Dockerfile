FROM golang:1.21.0-alpine AS builder
RUN apk add --no-cache gcc libc-dev
COPY . /src/
WORKDIR /src/cmd/shortlinks
RUN go build -mod=mod

FROM alpine:3.18
ENV TOKEN ""
COPY --from=builder /src/cmd/shortlinks/shortlinks /usr/local/bin/
COPY ./migrations /var/lib/shortlinks/migrations
COPY entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh"]
VOLUME "/data"
