# syntax=docker/dockerfile:1
# Create build stage based on buster image
FROM golang:1.20-alpine
# Create working directory under /app
WORKDIR /app
# Copy over all go config (go.mod, go.sum etc.)
COPY ./go.* ./
# Install any required modules
ENV GO111MODULE=on
RUN go mod init goapp && go mod tidy
# Copy over Go source code
COPY ./*.go ./
# Run the Go build and output binary under hello_go_http
RUN go build -o /goapp
# Run the app binary when we run the container
ENTRYPOINT ["/goapp"]