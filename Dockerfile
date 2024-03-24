ARG GO_VERSION=1.22.1
ARG ALPINE_VERSION=3.19

FROM golang:${GO_VERSION}-alpine AS builder

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY ./api/go.mod ./api/go.sum ./

RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY ./api/ .

RUN go build -o ./main ./cmd/main.go


ENV GIN_MODE release

CMD ["./main"]