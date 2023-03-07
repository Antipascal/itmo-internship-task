## Build
FROM golang:1.20 AS build

WORKDIR /

COPY ./ ./

#RUN go mod download
RUN go build -o /api

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /api /api

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/api"]