FROM golang:1.23 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /spacet ./cmd/spacet/main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /app 

COPY --from=build-stage /spacet /app/

CMD ["./spacet", "-config=/app/config/config.yaml"]