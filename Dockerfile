FROM golang AS builder

RUN mkdir /app
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest AS production

COPY --from=builder /app .
EXPOSE 8080
CMD ["./main"]
