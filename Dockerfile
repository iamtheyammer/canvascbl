FROM golang:alpine as build
COPY . /app
WORKDIR /app
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/canvasProxy src/main.go

FROM alpine
COPY --from=build /app/bin/canvasProxy /canvasProxy
EXPOSE 8000
ENTRYPOINT ["/canvasProxy"]

FROM nginx
