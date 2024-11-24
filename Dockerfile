FROM golang:1.23 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY data/ ./data/
COPY lib/ ./lib/
COPY types/ ./types/
COPY routes/ ./routes/
COPY templates/ ./templates/
COPY js/ ./js/
COPY css/ ./css/
COPY images/ ./images/

RUN CGO_ENABLED=1 GOOS=linux go build -o /server -a -ldflags '-linkmode external -extldflags "-static"' *.go

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM alpine AS build-release-stage

WORKDIR /

RUN apk --no-cache add --no-check-certificate ca-certificates \
  && update-ca-certificates

COPY --from=build-stage /server /server
COPY --from=build-stage /app/templates /templates

EXPOSE 3333

ENTRYPOINT ["/server"]
