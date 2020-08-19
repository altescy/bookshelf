FROM golang:1.14-alpine as builder
WORKDIR /build
COPY . ./
RUN CGO_ENABLED=off go build ./cmd/bookshelf


FROM alpine:latest
RUN apk --update add --no-cache tzdata libpq
WORKDIR /app
COPY --from=builder /build/bookshelf /app/bookshelf
CMD /app/bookshelf
