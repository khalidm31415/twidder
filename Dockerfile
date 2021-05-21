FROM golang:1.16-alpine AS build

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -o /bin/gin_twitter

FROM scratch
COPY --from=build /bin/gin_twitter /bin/gin_twitter
