FROM golang:alpine AS builder
RUN apk add --no-cache git

WORKDIR /app/

COPY . .
RUN go mod download

RUN CGO_ENABLED=0

RUN go build -o ./out/fantasyAPI ./cmd/APIServer/

FROM alpine:3.9 

COPY --from=builder /app/out/fantasyAPI /app/APIServer
COPY ./configs /configs

EXPOSE 8080

CMD ["/app/APIServer"]
