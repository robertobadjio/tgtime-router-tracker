FROM golang:1.22 as build

WORKDIR /app

COPY ../../go.mod go.sum ./
RUN go mod download

COPY ../../cmd/tracker/tracker.go ./
COPY ../../internal internal
COPY ../../config config
COPY ../../.env /

RUN CGO_ENABLED=0 GOOS=linux go build -o /tgtime-router-tracker

FROM scratch

COPY --from=build /tgtime-router-tracker .
COPY --from=ghcr.io/tarampampam/curl:8.6.0 /bin/curl /bin/curl
COPY --from=build /.env .

ENV EXPOSE_PORT=8081

CMD ["/tgtime-router-tracker"]