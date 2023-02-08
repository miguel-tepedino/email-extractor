FROM golang:latest
WORKDIR /api
COPY .env .env
COPY . .
RUN go mod download
RUN go build -o cmd/api cmd/api/apiv.go
EXPOSE 3002
CMD ["./cmd/api/apiv"]