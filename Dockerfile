FROM golang:1.22 as builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -o bin/notification-service .

FROM gcr.io/distroless/static-debian12
WORKDIR /
COPY --from=builder /workspace/bin/notification-service .

ENTRYPOINT ["/notification-service"]
