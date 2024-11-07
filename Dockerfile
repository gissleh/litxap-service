FROM docker.io/library/golang:1.23-alpine AS builder
WORKDIR /project
COPY . .
RUN go run ./cmd/litxap-fwewcheck
RUN go build ./cmd/litxap-service

FROM docker.io/library/alpine:latest
WORKDIR /project
COPY --from=builder /project/litxap-service /root/.fwew/dictionary-v2.txt ./
ENV PORT = 8080
CMD ./litxap-service

