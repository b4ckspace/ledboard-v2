FROM golang:1 AS build

WORKDIR /app
COPY . .

RUN go vet -v
RUN go test -v

RUN CGO_ENABLED=0 go build -o /tmp/ledboard

FROM gcr.io/distroless/static-debian12

COPY --from=build /tmp/ledboard /ledboard
CMD ["/ledboard"]
