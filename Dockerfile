FROM golang:latest AS builder
ENV GO111MODULE=on
ENV CGO_ENABLED=0
WORKDIR /workspace
COPY . /workspace
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o toto-server ./cmd/toto-server/main.go

FROM alpine
WORKDIR /app
RUN mkdir /migrations
COPY --from=builder /workspace/toto-server .
COPY --from=builder /workspace/config.yaml .
COPY --from=builder /workspace/migrations/* ./migrations/
CMD [ "./toto-server", "-config", "config.yaml" ]