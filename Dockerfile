ARG ZINCSEARCH_PASS=12345678
ARG ZINCSEARCH_USERNAME=lambda
FROM golang:latest
WORKDIR /api
COPY app.env app.env
COPY ./ .
RUN go mod download
RUN go build ./api
EXPOSE 3002
CMD ["./api/api"]