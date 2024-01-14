# Builder

FROM --platform=linux/amd64 golang:1.18-alpine as builder

WORKDIR /app

COPY . ./

RUN go build ./cmd/shipment-calculator-api

# Final docker image

FROM --platform=linux/amd64 alpine:3.7

WORKDIR /
COPY --from=builder /app/bin/shipment-calculator-api .

CMD ["/shipment-calculator-api"]
