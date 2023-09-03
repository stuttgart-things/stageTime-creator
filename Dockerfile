FROM golang:1.20 as builder
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

FROM eu.gcr.io/stuttgart-things/sthings-alpine@sha256:8d67f8b99f4bd4329cbdf5be80f8c8683a8a5cbe341c0860412c984b0a20a621
COPY --from=builder /bin/stageTime-creator /bin/stageTime-creator

ENTRYPOINT ["stageTime-creator"]
