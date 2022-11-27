# syntax=docker/dockerfile:1

# okay, let's make a volume mount and figure out how to mount a mysql instance to it, preferably the local-db-up instance

# also, it would be preferable to replace the base image w/ ubuntu @ some point, let's just do that now
# you know what, fuck it, I can play around with that shit later
FROM golang:1.16-alpine

ARG PROJECT_BINARY=bestir-permissionmaking-service
ARG PROJECT_BUILD_DIR=./build/bin

WORKDIR /app

RUN ls

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# these are previous failed attempts to copy subdirectories of the working directory of docker invocation into the working 
# directory of 
# COPY *.go ./cmd/bestir-permissionmaking-service
# COPY *.go ./internal/
# cmd/bestir-permissionmaking-service

# temporary measure
COPY *.go ./

RUN go build -o /bestir-permissionmaking-service

EXPOSE 8080

 CMD [ "/bestir-permissionmaking-service" ]

 