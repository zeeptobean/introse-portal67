FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
ARG SERVICE
RUN CGO_ENABLED=0 go build -trimpath -o /portal67 "./cmd/${SERVICE}"

FROM alpine:3.23
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /portal67 /usr/local/bin/portal67
COPY db ./db
USER 65534:65534
ENTRYPOINT ["portal67"]
