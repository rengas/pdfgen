FROM golang:1.18.2-buster AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY pkg pkg
COPY cmd/api cmd/api

RUN go build -o api ./cmd/api/

FROM gcr.io/distroless/base-debian10

WORKDIR /
COPY --from=build /app/api /bin/

EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT ["api"]
