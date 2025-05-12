FROM golang:1.24.3-alpine3.21 AS builder

RUN apk --no-cache add git protobuf protobuf-dev make

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
    && go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

ENV PATH="/go/bin:${PATH}"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make generateOpenapi
RUN make generateProto


RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/schedules ./cmd/schedules


FROM scratch

WORKDIR /app

COPY --from=builder /app/bin/schedules /app/schedules

COPY .env /app/.env
COPY db/migrations /app/migrations

ENTRYPOINT ["/app/schedules"]

EXPOSE 8080 12345