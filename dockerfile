FROM golang:alpine AS builder

RUN apk add --update --no-cache make

WORKDIR /app
COPY . .

RUN go mod download
RUN GOOS=linux go build -o ./bin/app ./cmd/main.go

FROM alpine:latest AS runner

COPY --from=builder /app/bin/app/ . 
COPY --from=builder /app/.env .
CMD [ "./app" ]