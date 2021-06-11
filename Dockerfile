FROM golang:1.16.4-alpine3.13 as builder

RUN adduser -D -g 'bytebot' bytebot
WORKDIR /app
COPY . .
RUN apk add --no-cache git tzdata
RUN ./docker-build.sh

FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/opt/bytebot /opt/bytebot
VOLUME /data

# Our chosen default for Prometheus
EXPOSE 8080
USER bytebot
ENTRYPOINT ["/opt/bytebot"]
