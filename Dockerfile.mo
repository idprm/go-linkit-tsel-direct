FROM golang:1.21-alpine as golang

RUN apk --no-cache add tzdata

RUN apk --update add ca-certificates

RUN mkdir -p /logs/http
RUN mkdir -p /logs/mo
RUN mkdir -p /logs/mt

WORKDIR /app
COPY . .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /linkit-tsel-direct .

FROM scratch

COPY --from=golang /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=golang /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=golang /etc/passwd /etc/passwd
COPY --from=golang /etc/group /etc/group
COPY --from=golang /bin/sh /bin/sh

COPY --from=golang /linkit-tsel-direct .

CMD ["/linkit-tsel-direct", "mo"]