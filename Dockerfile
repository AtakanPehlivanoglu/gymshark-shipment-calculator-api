#!/bin/sh
# Builder

FROM golang:1.18-alpine as builder

WORKDIR /app

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/shipment-calculator-api

# Final docker image

FROM alpine:3.7

ENV SPEC_FILE_PATH=config

WORKDIR /
COPY --from=builder /app/shipment-calculator-api .
COPY --from=builder /app .

CMD ["/shipment-calculator-api"]