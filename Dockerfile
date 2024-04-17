FROM golang:1.22.2-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN cd cmd/api && go build -o main main.go


# Running stage
FROM alpine:3.19

# Add this line to support health-check
RUN apk add curl

WORKDIR /app

COPY --from=builder /app/cmd/api/main .
COPY ./entrypoint.sh .

EXPOSE 8080

# ENTRYPOINT [ "/app/entrypoint.sh" ]
CMD [ "/app/main" ]
