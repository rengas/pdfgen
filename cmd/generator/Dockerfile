FROM golang:1.18.2-buster AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY pkg pkg
COPY cmd/generator cmd/generator

RUN go build -o generator ./cmd/generator/

FROM gcr.io/distroless/base-debian10

WORKDIR /
COPY --from=build /app/generator /bin/

EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT ["generator"]
