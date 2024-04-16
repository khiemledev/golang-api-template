FROM golang:1.22.2-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go


# Running stage
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/main .
COPY ./entrypoint.sh .

EXPOSE 8080

ENTRYPOINT [ "/app/entrypoint.sh" ]
CMD ["/app/main"]
