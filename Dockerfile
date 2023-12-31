FROM golang:1.21.4 AS builder
LABEL maintainer="Patrick Hermann patrick.hermann@sva.de"

ARG GO_MODULE="github.com/stuttgart-things/stageTime-creator"
ARG VERSION=""
ARG BUILD_DATE=""
ARG COMMIT=""

WORKDIR /src/
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /bin/stageTime-creator \
    -ldflags="-X ${GO_MODULE}/internal.version=${VERSION} -X ${GO_MODULE}/internal.date=${BUILD_DATE} -X ${GO_MODULE}/internal.commit=${COMMIT}"

RUN CGO_ENABLED=0 go build -o /bin/stcTestProducer tests/testProducer.go
RUN CGO_ENABLED=0 go build -o /bin/stcTestConsumer tests/testConsumer.go

FROM alpine:3.18.4
COPY --from=builder /bin/stageTime-creator /bin/stageTime-creator

# FOR SERVICE TESTING
COPY --from=builder /bin/stcTestProducer /bin/stcTestProducer
COPY --from=builder /bin/stcTestConsumer /bin/stcTestConsumer

ENTRYPOINT ["stageTime-creator"]
