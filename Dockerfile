# Compile stage
FROM golang:alpine AS build-env

# Add git to determine build git version
RUN apk add --no-cache --update git

# Add Delve
#RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Set GOPATH to build Go app
ENV GOPATH=/go

# Set apps source directory
ENV SRC_DIR=${GOPATH}/src/github.com/suaas21/godebug

# Define current working directory
WORKDIR ${SRC_DIR}

# Copy apps source code to the image
COPY . ${SRC_DIR}

# Build App
RUN ./build.sh


## Defining App image
FROM alpine:latest

RUN apk add --no-cache --update ca-certificates

WORKDIR /
#COPY --from=build-env /go/bin/dlv /
COPY --from=build-env /go/bin/godebug /
COPY config.yaml /

EXPOSE 8000 8005

#CMD ["/dlv", "--listen=:8005", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/godebug", "--", "serve", "-c", "/config.yaml"]
CMD ["/godebug", "serve", "-c", "/config.yaml"]