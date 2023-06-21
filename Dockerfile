FROM golang:alpine AS builder
WORKDIR /
COPY . .
# RUN CGO_ENABLED=0 go build -gcflags "-N -l" -o /app main.go
RUN go build -o /app main.go

FROM alpine:latest
COPY --from=builder /app /app
ENTRYPOINT ["/app"]