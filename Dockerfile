FROM golang:1.16-buster as build

WORKDIR /go/src/app
ADD . /go/src/app
RUN go build -o /go/bin/wrapdrive

FROM gcr.io/distroless/base-debian10
LABEL org.opencontainers.iamge.authors="secure@brickdoc.com"
COPY --from=build /go/bin/wrapdrive /
CMD ["/wrapdrive"]