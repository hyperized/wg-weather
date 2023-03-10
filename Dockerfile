FROM golang:latest as builder

WORKDIR /app/
COPY ./ .

RUN go mod tidy
RUN go test
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o weather .

FROM scratch

LABEL description="Weather in my local town as a Service"
LABEL author="Gerben Geijteman"

WORKDIR /app/

COPY --from=builder /app/weather /app/weather

EXPOSE 80/tcp
ENTRYPOINT ["/app/weather"]
CMD []
