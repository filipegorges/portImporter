FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /app/cmd/main /app/main

ENTRYPOINT ["/app/main"]
