FROM golang:1.20 as build-image
WORKDIR /go/src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o backend

FROM scratch
WORKDIR /root/
COPY --from=build-image /go/src/backend .
CMD ["./backend"]