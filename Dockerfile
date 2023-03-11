FROM alpine as certsbuidlder

RUN apk add -U --no-cache ca-certificates

FROM golang:latest as builder

WORKDIR /app/
COPY ./ .

RUN go mod tidy
RUN go test
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -buildvcs=false -o weather .

FROM scratch

LABEL description="Weather in my local town as a Service"
LABEL author="Gerben Geijteman"

WORKDIR /app/

COPY --from=builder /app/weather /app/weather
COPY --from=certsbuidlder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080/tcp
ENTRYPOINT ["/app/weather"]
CMD []
