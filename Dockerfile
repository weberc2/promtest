FROM golang:1.17.1 AS build

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY *.go .

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o app

FROM scratch

COPY --from=build /workspace/app /app

CMD ["/app"]
