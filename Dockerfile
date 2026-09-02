# syntax=docker/dockerfile:1

FROM golang:1.27 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/iris ./cmd/server

FROM gcr.io/distroless/base-debian12:nonroot
COPY --from=build /out/iris /iris

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/iris"]