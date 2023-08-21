FROM golang:1.20 as build-image
WORKDIR /go/src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o frontend

FROM scratch
WORKDIR /root/
COPY --from=build-image /go/src/frontend .
COPY ./static ./static
COPY ./templates ./templates
CMD ["./frontend"]