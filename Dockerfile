FROM golang:1.14-alpine as builder
WORKDIR /build
COPY . ./
RUN apk --update add --no-cache build-base
RUN CGO_ENABLED=off go build .


FROM alpine:latest
RUN apk --update add --no-cache tzdata
WORKDIR /app
COPY --from=builder /build/bookshelf /app/bookshelf
CMD /app/bookshelf
