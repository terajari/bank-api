FROM golang:1.21.3-alpine3.18 AS build
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine
WORKDIR /app
COPY --from=build /app/main .
COPY .env .

EXPOSE 8080
CMD [ "/app/main" ]


