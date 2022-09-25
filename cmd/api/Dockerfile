FROM golang:1.18.2-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY pkg pkg
COPY cmd/api cmd/api

RUN go build -o api ./cmd/api/


FROM alpine:3.16.2
FROM surnet/alpine-wkhtmltopdf:3.9-0.12.5-full

WORKDIR /
COPY --from=build /app/api /bin/
COPY --from=build /app/api/migrations /bin/

EXPOSE 8080


ENTRYPOINT ["api"]
