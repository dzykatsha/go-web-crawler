FROM golang:1.19-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /go/bin/api ./cmd/api/main.go 

FROM gcr.io/distroless/base-debian11

COPY --from=build /go/bin/api /

EXPOSE 8000

ENTRYPOINT ["/api"]
