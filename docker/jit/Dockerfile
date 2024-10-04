FROM golang:1.21.3-alpine
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build ./cmd/query_jit_microservice
EXPOSE 8001
CMD ["./query_jit_microservice"]