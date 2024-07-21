FROM golang:1.22 as builder
ENV CGO_ENABLED=0
COPY . /
WORKDIR /
RUN go build -o cmd/finance .

FROM scratch
COPY --from=builder /cmd/finance /finance
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /
CMD ["./finance"] 