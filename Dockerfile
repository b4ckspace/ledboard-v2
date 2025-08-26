FROM golang:1 AS build

WORKDIR /app
COPY . .

RUN go mod download
RUN go vet -v
RUN go test -v

RUN CGO_ENABLED=0 go build -o /app/ledboard

FROM gcr.io/distroless/static-debian12

COPY --from=build /app/ledboard /
CMD ["/app/ledboard"]
