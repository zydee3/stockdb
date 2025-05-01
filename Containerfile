FROM golang:1.24 AS build

WORKDIR /work

COPY . ./

RUN make build

FROM build AS test-stage
RUN make test

## Deploy
FROM ubuntu:latest

WORKDIR /

COPY --from=build /work/build/stockd /bin/stockd
COPY --from=build /work/build/stockctl /bin/stockctl

ENTRYPOINT ["/bin/stockd"]
