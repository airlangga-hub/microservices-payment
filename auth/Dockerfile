FROM golang:latest AS BUILD

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o auth ./main.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/auth .

ENTRYPOINT [ "./auth" ]