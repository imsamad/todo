# Build stage
FROM golang:1.22 as build
WORKDIR /app
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download
COPY . .
RUN go build -o /build/app .

# Run stage
FROM golang:1.22 as run
WORKDIR /app
COPY --from=build /build/app /app
EXPOSE 4000
ENTRYPOINT [ "/app" ]
