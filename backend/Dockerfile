FROM golang:latest AS build

WORKDIR /app

# Copy the Go module files
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o cmd ./src

FROM busybox:latest AS run

COPY --from=build /app/cmd /cmd

EXPOSE 8080
CMD ["/cmd"]
