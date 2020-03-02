# build stage
FROM golang:1.13.4 AS builder

ENV GO111MODULE auto
ENV CGO_ENABLED=0

ADD . /src
WORKDIR /src
RUN make build

# final stage
FROM scratch

WORKDIR /app
COPY --from=builder /src/bin/ecr-lifecycle /app/
ENTRYPOINT ["./ecr-lifecycle"]
